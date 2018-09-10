package session

import (
    "log"
    "github.com/astaxie/beego/config"
    "encoding/base64"
    "encoding/json"
    "net/url"
    "crypto/aes"
    "crypto/cipher"
    "github.com/wulijun/go-php-serialize/phpserialize"
    "io/ioutil"
    "errors"
)

var appKey []byte
var host string
var projectPath string

func init() {
    iniconf, err := config.NewConfig("ini", "conf/laravel.conf")
    if err != nil {
        log.Fatal(err)
    }
    projectPath = iniconf.String("projectpath")

    envconf, err := config.NewConfig("ini", projectPath + "/.env")
    if err != nil {
        log.Fatal(err)
    }

    key := envconf.String("APP_KEY")
    if len(key) <51 {
        log.Fatal("Значение APP_KEY неверно или не задано в .env файле")
    }
    key = key[7:]
    appKey, err = base64.StdEncoding.DecodeString(key)

    if err != nil {
            log.Fatal(err)
    }
}


func GetUserId(cookie string) (int, error) {
    str, err := url.PathUnescape(cookie)
    if err != nil {
            return -1, err
    }

    data, err := base64.StdEncoding.DecodeString(str)
    if err != nil {
            return -1, err
    }

    var dat map[string][]byte
    if err := json.Unmarshal(data, &dat); err != nil {
        return -1, err
    }

    block, err := aes.NewCipher(appKey)
    if err != nil {
        log.Print(err)
        return -1, err
    }

    cbc := cipher.NewCBCDecrypter(block, dat["iv"])
    cbc.CryptBlocks(dat["value"], dat["value"])

    cleartext, err := phpserialize.Decode(string(dat["value"]))
    if err != nil {
        return -1, err
    }
    sessionFile := projectPath + "/storage/framework/sessions/" + cleartext.(string)
    buf, err := ioutil.ReadFile(sessionFile)
    if err != nil {
        return -1, err
    }

    cleartext, err = phpserialize.Decode(string(buf))
    if err != nil {
        return -1, err
    }

    decodeData, ok := cleartext.(map[interface{}]interface{})
    if !ok {
        err = errors.New("Cant convert DecodeData...")
        return -1, err
    }
    // log.Print(decodeData)
    result, _ := decodeData["login_web_59ba36addc2b2f9401580f014c7f58ea4e30989d"].(int64)

    return int(result), nil
}
