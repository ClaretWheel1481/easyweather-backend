package main

import (
	"EasyWeatherBackend/auth"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	// TODO: 请输入你的高德与和风api key
	var gaodeKey string = "please input a key"
	var qWeatherKey1 string = "please input a key"
	var qWeatherKey2 string = "please input a key"

	gin.SetMode(gin.ReleaseMode)

	// 路由组
	ewRouter := router.Group("/v1/data")
	ewRouter.Use(auth.AuthMiddleware())
	{
		// 查询市、区id
		ewRouter.GET("/baseCityInfo/:cityName", func(c *gin.Context) {
			cityName := c.Param("cityName")
			url := fmt.Sprintf("https://restapi.amap.com/v3/config/district?keywords=%s&subdistrict=0&key=%s&extensions=base", cityName, gaodeKey)
			GetUrlJsonData(url, c)
		})

		// 查询基础天气情况
		ewRouter.GET("/baseWeatherInfo/:cityId", func(c *gin.Context) {
			cityId := c.Param("cityId")
			url := fmt.Sprintf("https://restapi.amap.com/v3/weather/weatherInfo?city=%s&key=%s&extensions=base", cityId, gaodeKey)
			GetUrlJsonData(url, c)
		})

		// 查询全部天气情况
		ewRouter.GET("/allWeatherInfo/:cityId", func(c *gin.Context) {
			cityId := c.Param("cityId")
			url := fmt.Sprintf("https://restapi.amap.com/v3/weather/weatherInfo?city=%s&key=%s&extensions=all", cityId, gaodeKey)
			GetUrlJsonData(url, c)
		})

		// 获取和风天气的cityCode
		ewRouter.GET("/CityId/:adCode", func(c *gin.Context) {
			adCode := c.Param("adCode")
			url := fmt.Sprintf("https://geoapi.qweather.com/v2/city/lookup?location=%s&key=%s", adCode, qWeatherKey1)
			GetUrlJsonData(url, c)
		})

		// 获取天气预警
		ewRouter.GET("/CityWarning/:qWeatherCityId", func(c *gin.Context) {
			qWeatherCityId := c.Param("qWeatherCityId")
			url := fmt.Sprintf("https://devapi.qweather.com/v7/warning/now?location=%s&lang=cn&key=%s", qWeatherCityId, qWeatherKey1)
			GetUrlJsonData(url, c)
		})

		// 获取天气指数
		ewRouter.GET("/CityIndices/:qWeatherCityId", func(c *gin.Context) {
			qWeatherCityId := c.Param("qWeatherCityId")
			url := fmt.Sprintf("https://devapi.qweather.com/v7/indices/1d?type=1,2&location=%s&key=%s", qWeatherCityId, qWeatherKey2)
			GetUrlJsonData(url, c)
		})

		// 获取空气质量
		ewRouter.GET("/CityAir/:qWeatherCityId", func(c *gin.Context) {
			qWeatherCityId := c.Param("qWeatherCityId")
			url := fmt.Sprintf("https://devapi.qweather.com/v7/air/now?location=%s&key=%s", qWeatherCityId, qWeatherKey2)
			GetUrlJsonData(url, c)
		})
	}
	// 生成令牌的路由
	router.GET("/generateToken", func(c *gin.Context) {
		token, err := auth.GenerateToken()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"token": token})
	})
	router.Run(":37878")
}

// GetUrlJsonData 获取URL的JSON数据
func GetUrlJsonData(url string, c *gin.Context) {
	resp, err := http.Get(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resp.Body)
	var data interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}
