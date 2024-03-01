package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	router := gin.Default()

	ewRouter := router.Group("/v1/data")
	{
		//查询市、区id
		ewRouter.GET("/baseCityInfo/:cityName", func(c *gin.Context) {
			cityName := c.Param("cityName")
			url := fmt.Sprintf("https://restapi.amap.com/v3/config/district?keywords=%s&subdistrict=0&key=1c918e3c95e677120dc2a54bf58a408f&extensions=base", cityName)
			getUrlJsonData(url, c)
		})

		//查询基础天气情况
		ewRouter.GET("/baseWeatherInfo/:cityId", func(c *gin.Context) {
			cityId := c.Param("cityId")
			url := fmt.Sprintf("https://restapi.amap.com/v3/weather/weatherInfo?city=%s&key=1c918e3c95e677120dc2a54bf58a408f&extensions=base", cityId)
			getUrlJsonData(url, c)
		})

		//查询全部天气情况
		ewRouter.GET("/allWeatherInfo/:cityId", func(c *gin.Context) {
			cityId := c.Param("cityId")
			url := fmt.Sprintf("https://restapi.amap.com/v3/weather/weatherInfo?city=%s&key=1c918e3c95e677120dc2a54bf58a408f&extensions=all", cityId)
			getUrlJsonData(url, c)
		})

		//根据高德开放平台的adCode获取和风天气的cityCode
		ewRouter.GET("/getCityId/:adCode", func(c *gin.Context) {
			adCode := c.Param("adCode")
			url := fmt.Sprintf("https://geoapi.qweather.com/v2/city/lookup?location=%s&key=aefbbb5f60264c0fbb97a5ecca6ce2d3", adCode)
			getUrlJsonData(url, c)
		})

		//根据和风天气cityCode获取天气预警
		ewRouter.GET("/getCityWarning/:qWeatherCityId", func(c *gin.Context) {
			qWeatherCityId := c.Param("qWeatherCityId")
			url := fmt.Sprintf("https://devapi.qweather.com/v7/warning/now?location=%s&lang=cn&key=aefbbb5f60264c0fbb97a5ecca6ce2d3", qWeatherCityId)
			getUrlJsonData(url, c)
		})

		//根据和风天气cityCode获取天气指数
		ewRouter.GET("/getCityIndices/:qWeatherCityId", func(c *gin.Context) {
			qWeatherCityId := c.Param("qWeatherCityId")
			url := fmt.Sprintf("https://devapi.qweather.com/v7/indices/1d?type=1,2&location=%s&key=aefbbb5f60264c0fbb97a5ecca6ce2d3", qWeatherCityId)
			getUrlJsonData(url, c)
		})

		//根据和风天气cityCode获取空气质量
		ewRouter.GET("/getCityAir/:qWeatherCityId", func(c *gin.Context) {
			qWeatherCityId := c.Param("qWeatherCityId")
			url := fmt.Sprintf("https://devapi.qweather.com/v7/air/now?location=%s&key=aefbbb5f60264c0fbb97a5ecca6ce2d3", qWeatherCityId)
			getUrlJsonData(url, c)
		})
	}

	router.Run(":37878")
}

func getUrlJsonData(url string, c *gin.Context) {
	resp, err := http.Get(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()
	var data interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}
