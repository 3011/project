# portScanner-go


因为没找到特别顺手的端口扫描工具，所以造了个轮子。

并发扫描，应该能在十秒内扫描所有端口。

（linux 与 macOS 版本未测试也许有问题）

## 示例（以 Windows 版本为例）：
通过 TCP 与 UDP 扫描 example.com 所有端口:
```cmd
.\portScanner-go-amd64.exe example.com
```


通过 TCP 与 UDP 扫描 127.0.0.1 所有端口:
```cmd
.\portScanner-go-amd64.exe 127.0.0.1
```


通过 TCP  扫描 127.0.0.1 的 80 至 5000 端口:
```cmd
.\portScanner-go-amd64.exe -tcp -port 80,5000 127.0.0.1
```


通过 UDP  扫描 127.0.0.1 的 50 至 50000 端口:
```cmd
.\portScanner-go-amd64.exe -tcp -port 50,50000 127.0.0.1
```


**注意：域名或IP需放在参数后**
错误示范：
```cmd
.\portScanner-go-amd64.exe 127.0.0.1 -tcp -port 50,50000
```