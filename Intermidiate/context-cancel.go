package main

import (
	"context"
	"fmt"
	"time"
	"net/http"
)
func fetchdata(ctx context.Context)error{
	req,err:=http.NewRequestWithContext(ctx,"GET","https://httpbin.org/delay/1",nil)
	if err!=nil{
		return err}
	client:=&http.Client{}
	resp,err:=client.Do(req)	
	if err!=nil{
		return err}
	defer resp.Body.Close()
	fmt.Println("Data fetched successfully",resp.StatusCode)
	return nil
}

func main(){
	ctx,cancel:=context.WithTimeout(context.Background(),3*time.Second)
	defer cancel()
	err:=fetchdata(ctx)
	if err!=nil{
		fmt.Println("Error fetching data:",err)
	
		return
	}
		fmt.Println("Data fetched successfully")	


}