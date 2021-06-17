## eth
```go
package ...
import (
    ethSdk "github.com/sunliang711/eth/sdk"
)

func main(){
	...
	txManager, err := ethSdk.New(rcp,price,limit,timeout,interval)
	if err != nil {
		...
    }
    txManager.WriteContractSync(...)
	...
}

```