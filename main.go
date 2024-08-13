package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("LanceHuang")

func main() {
	router := gin.Default()
	var gaodeKey string = "1c918e3c95e677120dc2a54bf58a408f"
	var qWeatherKey1 string = "d5ae8d9521c148879fc4268e6fcf624e"
	var qWeatherKey2 string = "325e370be5254274a7460665d95559b9"

	gin.SetMode(gin.ReleaseMode)

	// 路由组
	ewRouter := router.Group("/v1/data")
	ewRouter.Use(AuthMiddleware())
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
		token, err := GenerateToken()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"token": token})
	})
	router.Run(":37878")
}

// AuthMiddleware 验证令牌
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return secretKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// GenerateToken 生成新的JWT令牌
func GenerateToken() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})
	return token.SignedString(secretKey)
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
