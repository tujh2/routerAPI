package main

import (
    jwt "github.com/appleboy/gin-jwt"
    "github.com/gin-gonic/gin"
    "log"
    "net/http"
    "time"
)

var configPath = "/etc/"
var configName = "routerAPI.conf"
var secretKey = "safklweh8ry3njdspf032"
var identityKey = "id"

func main() {
    readConfig()

    r := gin.New()
    r.Use(gin.Logger())
    r.Use(gin.Recovery())

    authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
        Realm:       "test zone",
        Key:         []byte(secretKey),
        Timeout:     time.Hour,
        MaxRefresh:  time.Hour,
        IdentityKey: identityKey,
        LoginResponse: func(c *gin.Context, code int, token string, expire time.Time) {
            c.JSON(http.StatusOK, gin.H{
                "token":  token,
                "expire": expire.Format(time.RFC3339),
            })
        },

        PayloadFunc: func(data interface{}) jwt.MapClaims {
            if v, ok := data.(*JSONLogin); ok {
                var user = jwt.MapClaims{
                    identityKey: v.Username,
                }
                return user
            }
            return jwt.MapClaims{}
        },
        IdentityHandler: func(c *gin.Context) interface{} {
            claims := jwt.ExtractClaims(c)
            return &JSONLogin{
                Username: claims[identityKey].(string),
            }
        },
        Authenticator: func(c *gin.Context) (interface{}, error) {
            var login JSONLogin
            if err := c.ShouldBind(&login); err != nil {
                return "", jwt.ErrMissingLoginValues
            }
            if login == adminAuthUser {
                return &adminAuthUser, nil
            }
            return nil, jwt.ErrFailedAuthentication
        },
        Authorizator: func(data interface{}, c *gin.Context) bool {
            return true
        },
        Unauthorized: func(c *gin.Context, code int, message string) {
            c.JSON(code, gin.H{
                "code":    code,
                "message": message,
            })
        },

        TokenLookup: "header: Authorization, query: token, cookie: jwt",

        TokenHeadName: "Bearer",

        TimeFunc: time.Now,
    })

    if err != nil {
        log.Fatal("JWT Error:" + err.Error())
    }

    r.POST("/login", authMiddleware.LoginHandler)

    r.NoRoute(authMiddleware.MiddlewareFunc(), func(c *gin.Context) {
        claims := jwt.ExtractClaims(c)
        log.Printf("NoRoute claims: %#v\n", claims)
        c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
    })

    auth := r.Group("/auth")

    // Refresh time can be longer than token timeout
    auth.GET("/refresh_token", authMiddleware.RefreshHandler)
    auth.Use(authMiddleware.MiddlewareFunc())
    {
    }

    if err := http.ListenAndServe(IP+":"+PORT, r); err != nil {
        log.Fatal(err)
    }
}
