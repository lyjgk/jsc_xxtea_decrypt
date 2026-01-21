package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	flag "github.com/spf13/pflag"
	"github.com/xxtea/xxtea-go/xxtea"
)

func example() {
	str := "Hello World! 你好，中国！"
	key := "1234567890"
	encrypt_data := xxtea.Encrypt([]byte(str), []byte(key))
	decrypt_data := string(xxtea.Decrypt(encrypt_data, []byte(key)))
	if str == decrypt_data {
		fmt.Println("success!")
	} else {
		fmt.Println("fail!")
	}
}

func decryptJSC(jscPath, key string) error {
	data, err := os.ReadFile(jscPath)
	if err != nil {
		return err
	}
	decrypt := xxtea.Decrypt(data, []byte(key))
	reader, err := gzip.NewReader(bytes.NewReader(decrypt))
	if err != nil {
		return err
	}
	defer reader.Close()

	unzipped, err := io.ReadAll(reader)
	if err != nil {
		return err
	}

	newPath := strings.TrimSuffix(jscPath, filepath.Ext(jscPath)) + ".js"
	return os.WriteFile(newPath, unzipped, 0644)
}

func main() {
	// 定义两个 string 参数
	path := flag.String("path", "./jsc/", "文件路径")
	key := flag.String("key", "", "密钥")

	// 解析命令行参数
	flag.Parse()
	if *key == "" {
		fmt.Println("请输入密钥")
		return
	}
	// 使用参数
	fmt.Println("path:", *path)
	fmt.Println("key:", *key)

	filepath.WalkDir(*path, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		if filepath.Ext(path) == ".jsc" {
			if err := decryptJSC(path, *key); err != nil {
				fmt.Printf("filename: %s err: %s\n", path, err)
			} else {
				fmt.Printf("filename: %s success\n", path)
			}
		}
		return nil
	})
}
