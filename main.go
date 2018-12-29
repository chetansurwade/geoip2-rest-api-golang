package main

import (
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gin-contrib/cache"
	"github.com/gin-contrib/cache/persistence"
	"github.com/gin-gonic/gin"

	geoip2 "github.com/oschwald/geoip2-golang"
)

// UserIP for user's ip info
type UserIP struct {
	UserIP []string `form:"ip" json:"ip" xml:"ip"  binding:"required"`
}

// Cors for cross origin resource sharing in header
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		if c.Request.Method != "OPTIONS" {
			c.Next()
		} else {
			c.Writer.Header().Add("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
			c.Writer.Header().Add("Access-Control-Allow-Headers", "authorization, origin, content-type, accept")
			c.Writer.Header().Add("Allow", "HEAD,GET,POST,PUT,PATCH,DELETE,OPTIONS")
			c.Writer.Header().Add("Content-Type", "application/json")
			c.AbortWithStatus(http.StatusOK)
		}
	}
}

// main function containing the routes and db initialization
func main() {
	// set gin to production mode
	gin.SetMode(gin.ReleaseMode)

	var citydb, err = geoip2.Open("GeoLite2-City.mmdb")
	if err != nil {
		log.Fatal(err)
	}

	defer citydb.Close()
	r := gin.Default()
	r.Use(Cors())
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "API for GeoIP details",
			"status":  true,
		})
	})
	// for faster response an inmemory cache store of routes
	store := persistence.NewInMemoryStore(time.Second)
	// Supports any method GET/POST along with any content-type like json,form xml
	// change the time.Minute to your choice of duration
	// supports a single ip as input, for example http://localhost:8080/geoip?ip=YOUR-IP or http://localhost:8080/geoip?ip=YOUR-IP&ip=YOUR-IP2&ip=YOUR-IP3
	r.Any("/geoip", cache.CachePage(store, time.Minute, func(c *gin.Context) {
		var usrIP UserIP
		err := c.ShouldBind(&usrIP)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":  "bad request",
				"status": false,
			})
			return
		}
		userIPs := usrIP.UserIP
		if len(userIPs) < 1 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":  "Kindly specify the ip or array of ips",
				"status": false,
			})
			return
		}
		var results []interface{}

		for _, userIP := range userIPs {
			data := make(map[string]interface{}, 0)
			data["ip"] = userIP
			ip := net.ParseIP(userIP)
			if ip == nil {
				data["error"] = true
				results = append(results, data)
				continue
			}

			cityRecord, err := citydb.City(ip)
			if err != nil {
				data["error"] = true
				results = append(results, data)
				continue
			} else {
				var division, divisionCode string
				if len(cityRecord.Subdivisions) > 0 {
					division = cityRecord.Subdivisions[0].Names["en"]
					divisionCode = cityRecord.Subdivisions[0].IsoCode
				}
				data["continent"] = cityRecord.Continent.Names["en"]
				data["continent_code"] = cityRecord.Continent.Code
				data["country"] = cityRecord.Country.Names["en"]
				data["country_code"] = cityRecord.Country.IsoCode
				data["error"] = false
				data["division"] = division
				data["division_code"] = divisionCode
				data["city"] = cityRecord.City.Names["en"]
				data["postal_code"] = cityRecord.Postal.Code
				data["lat"] = cityRecord.Location.Latitude
				data["lon"] = cityRecord.Location.Longitude
				data["accuracy"] = cityRecord.Location.AccuracyRadius
				data["timezone"] = cityRecord.Location.TimeZone
				results = append(results, data)
			}
		}
		c.JSON(http.StatusOK, gin.H{
			"result": results,
			"status": true,
		})
		return
	}))

	r.Run(":8080")
	return
}
