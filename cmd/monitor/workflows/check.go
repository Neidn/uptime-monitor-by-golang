package workflows

import (
	"fmt"
	"github.com/Neidn/uptime-monitor-by-golang/config"
	"github.com/gorilla/websocket"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var FailedPerformanceTest = PerformanceTestResult{
	Status:       StatusDown,
	ResponseTime: 0,
	Result: PerformanceTestCode{
		HttpCode: 0,
	},
}

func TcpPingCheck(site config.Site) (PerformanceTestResult, error) {
	log.Println("Using TCP Ping to check site")
	var responseTime time.Duration
	var maxResponseTime time.Duration
	if site.MaxResponseTime > 0 {
		maxResponseTime = time.Duration(site.MaxResponseTime) * time.Millisecond
	} else {
		maxResponseTime = time.Duration(config.MaxResponseTime) * time.Millisecond
	}

	address := site.Url
	status := StatusUp

	ip := net.ParseIP(site.Url)

	if ip == nil {
		_url, err := url.Parse(site.Url)
		if err != nil {
			log.Println("Error parsing URL", err)
			return FailedPerformanceTest, err
		}
		ipList, err := net.LookupIP(fmt.Sprintf("%s", _url))

		if err != nil {
			log.Println("Error looking up IP", err)
			return PerformanceTestResult{}, err
		}

		address = ipList[0].String()

	}

	// with specific port
	if site.Port != 0 {
		// check if the address is ipv6 or ipv4
		if strings.Contains(address, ":") {
			address = fmt.Sprintf("[%s]:%d", address, site.Port)
		} else {
			address = fmt.Sprintf("%s:%d", address, site.Port)
		}
	}

	var responseTimeList []time.Duration

	// ping the address
	// Repeat the ping for the default count
	for i := 0; i < config.TcpPingDefaultCount; i++ {
		startTime := time.Now()
		conn, err := net.DialTimeout("tcp", address, time.Duration(config.TcpPingDefaultTimeout)*time.Second)

		if err != nil {
			log.Println("Error dialing", err)
			status = StatusDown
			continue
		}
		defer func(conn net.Conn) {
			err := conn.Close()
			if err != nil {
				log.Println("Error closing connection", err)
			}
		}(conn)

		responseTimeList = append(responseTimeList, time.Since(startTime))

		time.Sleep(time.Duration(config.TcpPingDefaultInterval) * time.Second)
	}

	responseTime = avg(responseTimeList)

	if responseTime < 0 {
		return FailedPerformanceTest, nil
	}

	if responseTime > maxResponseTime {
		status = StatusDegraded
	}

	return PerformanceTestResult{
		Result: PerformanceTestCode{
			HttpCode: 200,
		},
		ResponseTime: int(responseTime.Milliseconds()),
		Status:       status,
	}, nil
}

func WsCheck(site config.Site) (PerformanceTestResult, error) {
	log.Println("Using WebSocket to check site")
	status := StatusUp

	c, _, err := websocket.DefaultDialer.Dial(site.Url, nil)
	if err != nil {
		log.Println("Error dialing", err)
		return FailedPerformanceTest, err
	}
	defer func(c *websocket.Conn) {
		err := c.Close()
		if err != nil {
			log.Println("Error closing connection", err)
		}
		log.Println("WebSocket connection closed")
	}(c)

	_textMessage := []byte("")
	if site.Body != "" {
		_textMessage = []byte(site.Body)
	}

	err = c.WriteMessage(websocket.TextMessage, _textMessage)
	if err != nil {
		status = StatusDown
	}

	return PerformanceTestResult{
		Result: PerformanceTestCode{
			HttpCode: 200,
		},
		Status:       status,
		ResponseTime: 0,
	}, nil
}

var DefaultExpectedStatusCodes = map[int]bool{
	200: true,
	201: true,
	202: true,
	203: true,
	204: true,
	205: true,
	206: true,
	207: true,
	208: true,
	226: true,
	300: true,
	301: true,
	302: true,
	303: true,
	304: true,
	305: true,
	306: true,
	307: true,
	308: true,
}

func HttpCheck(site config.Site) (PerformanceTestResult, error) {
	log.Println("Using HTTP to check site")
	var expectedStatusCodes map[int]bool
	var responseTime time.Duration
	var maxResponseTime time.Duration
	if site.MaxResponseTime > 0 {
		maxResponseTime = time.Duration(site.MaxResponseTime) * time.Millisecond
	} else {
		maxResponseTime = time.Duration(config.MaxResponseTime) * time.Millisecond
	}
	status := StatusUp
	startTime := time.Now()

	if site.ExpectedStatusCode > 0 {
		expectedStatusCodes = map[int]bool{
			site.ExpectedStatusCode: true,
		}
	} else {
		expectedStatusCodes = DefaultExpectedStatusCodes
	}

	resp, err := http.Get(site.Url)
	if err != nil {
		log.Println("Error getting response", err)
		return FailedPerformanceTest, err
	}

	data := make([]byte, 1024)
	_, err = resp.Body.Read(data)
	if err != nil {
		log.Println("Error reading response body", err)
		return FailedPerformanceTest, err
	}

	bodyData := string(data)

	defer func(resp *http.Response) {
		err := resp.Body.Close()
		if err != nil {
			log.Println("Error closing response body", err)
		}
	}(resp)

	responseTime = time.Since(startTime)

	if !expectedStatusCodes[resp.StatusCode] {
		status = StatusDown
	}

	if responseTime > maxResponseTime {
		status = StatusDegraded
	}

	if status == StatusUp {
		if site.DangerousBodyDown != "" && strings.Contains(bodyData, site.DangerousBodyDown) {
			status = StatusDown
		}

		if site.DangerousBodyDegraded != "" && strings.Contains(bodyData, site.DangerousBodyDegraded) {
			status = StatusDegraded
		}
	}

	if site.DangerousBodyDegradedIfTextMissing != "" && !strings.Contains(bodyData, site.DangerousBodyDegradedIfTextMissing) {
		status = StatusDegraded
	}

	if site.DangerousBodyDownIfTextMissing != "" && !strings.Contains(bodyData, site.DangerousBodyDownIfTextMissing) {
		status = StatusDown
	}

	return PerformanceTestResult{
		Result: PerformanceTestCode{
			HttpCode: resp.StatusCode,
		},
		Status:       status,
		ResponseTime: int(responseTime.Milliseconds()),
	}, nil
}

func avg(responseTimeList []time.Duration) time.Duration {
	var sum time.Duration
	for _, v := range responseTimeList {
		sum += v
	}
	return sum / time.Duration(len(responseTimeList))
}
