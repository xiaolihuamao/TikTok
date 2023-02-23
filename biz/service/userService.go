package service

import (
	"TikTok/biz/dao"
	"TikTok/biz/model"
	"TikTok/biz/service/mysql"
	"TikTok/conf"
	"bytes"
	"context"
	"crypto/cipher"
	"crypto/des"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"unsafe"
)

func GetuserInfo(ctx context.Context, c *app.RequestContext, id int64) ([]model.User, error) {
	var userList []model.User

	dao.Db.Table("users").Select("user_id", "username", "follow_count", "follower_count").Distinct().Where("user_id = ?", id).Scan(&userList)
	return userList, nil
}

// 将新注册的用户,插入到数据库中
func Registeruser(ctx context.Context, c *app.RequestContext, username string, password string) (userid int64, err error) {

	usernames := make([]string, 0)
	dao.Db.Table("users").Select("username").Distinct().Scan(&usernames) //Scan与Find
	if IsInSlice(username, usernames) {                                  //如果有的话那么则返回用户名已经存在
		return 0, errors.New("用户名已经存在")
	}
	key1 := []byte(conf.Secret) //密钥进行匹配
	result1 := DesCbcEncryption([]byte(password), key1)
	fmt.Println(result1)
	passwordaes := base64.StdEncoding.EncodeToString(result1)
	fmt.Println(passwordaes)
	//reward1 := desCbcDecryption(result1, key1)
	//fmt.Printf("%s\n", reward1)
	user := model.User{Username: username, Password: passwordaes, FollowCount: 0, FollowerCount: 0}
	//用户的列表，，用户名，密码，关注列表，被关注列表
	/*dao.Db.Create(&user)
	userid, err = user.UserID, result.Error // 返回插入记录的条数// 返回 error
	return userid, err*/
	return mysql.InsertUser(&user)
}

// 判断一个元素是否在一个数组中
func IsInSlice(element string, elements []string) (isIn bool) {
	for _, item := range elements {
		if element == item {
			isIn = true
			return
		}
	}
	return
}

//加密模块
/*
工具方法：
(1)工具函数，实现对最后一个分组数据进行填充
	supplementLastGroup
	plainText []byte : 明文切片
	blockSize int : 根据算法决定分组的长度 （des、3des是8，aes是16）
(2)工具函数，实现对最后一个分组数据进行填充取消
	unsupplementLastGroup
	plainText []byte : 明文切片
*/
func supplementLastGroup(plainText []byte, blockSize int) []byte {
	//获取填充长度
	supNum := blockSize - len(plainText)%blockSize
	//创建一个新的字符切片，长度和填充长度相同，用填充长度填充（却几补几）
	supText := bytes.Repeat([]byte{byte(supNum)}, supNum)
	//把supText拼接到plainText内部，然后返回
	//...可以监测到切片，并对切片进行扩容
	return append(plainText, supText...)
}
func unsupplementLastGroup(plainText []byte) []byte {
	//获取最后一个字符，并还原回整数
	lastCharNum := int(plainText[len(plainText)-1])
	//将原内容从起始位置截取到【长度-补充】位置。
	return plainText[:len(plainText)-lastCharNum]
}

// 使用【DES加密算法】+【CBC分组密码模式】加密--解密
func DesCbcEncryption(plainText, key []byte) []byte {
	//1.创建一个底层加密接口对象（des，3des还是aes）
	block, err := des.NewCipher(key)
	if err != nil {
		panic(err)
	}
	//2.如果选择ecb或者cbc密码分组模式，则需要对明文最后一个分组内容进行填充
	newText := supplementLastGroup(plainText, block.BlockSize()) //supplementLastGroup(plainText,des.BlockSize)
	//3.创建一个密码分组模式的接口对象（cbc和ctr）
	iv := []byte{1, 2, 3, 4, 5, 6, 7, 8}
	blcokMode := cipher.NewCBCEncrypter(block, iv)
	//4.实现加密
	blcokMode.CryptBlocks(newText, newText)
	//5.返回加密密文
	return newText
}
func desCbcDecryption(cipherText, key []byte) []byte {
	//1.创建一个底层加密接口对象（des，3des还是aes）
	block, err := des.NewCipher(key)
	if err != nil {
		panic(err)
	}
	//2.创建一个密码分组模式的接口对象（cbc和ctr）
	iv := []byte{1, 2, 3, 4, 5, 6, 7, 8}
	blcokMode := cipher.NewCBCDecrypter(block, iv)
	//3.实现解密
	blcokMode.CryptBlocks(cipherText, cipherText)
	//4.对最后一个分组数据进行消除填充,并返回
	return unsupplementLastGroup(cipherText)
}

// []byte转string
func Sbyte2str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
