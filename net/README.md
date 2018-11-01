# goutils for net utils:

## history

```
GetLocalNS() (NSRecords, error)      
//get local dns config list
```

```
func GetGatewayByNic(nicName string) (net.IP, error)    
//return gateway by nicName
```

```
func ListGateway() ([]*RouteItem, error) 
//list all gateway of each nic
```

```
func GetDefaultGateway() (net.IP, error) 
```

```
func IsPublicIP(IP net.IP) bool 
//if an ip is public ip
```

```
func GetHostIP() (string, error) 
//get local host ip, like hostname -I
```

```
func GetInternalIP() (string, error) {
// get the first ip of nic (will skip loopback)
```

```
func GetRequestIP(r *http.Request) string 
// get Request clientIP in http request
```

```
func GetInternalIPByDevName(dev string) ([]string, error) {
// get nic ip list by name	
```	

```
func GetUploadFileSize(upfile multipart.File) (int64, error) {
// get http upload file size
```

```
func GetLocalConnectIP(proto string, addr string) (string, error) {
//
```