apple.pem //苹果公共证书
private.pem //苹果私钥
config.yaml
    kid: 
    iss: 
    bid: 

```
生成bundled.go

 fyne bundle private.pem > bundled.go
 fyne bundle -append apple.pem >> bundled.go
 fyne bundle -append config.yaml >> bundled.go
 fyne bundle -append simkai.ttf >> bundled.go
```