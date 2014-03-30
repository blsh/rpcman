package rpcman

import (
    "fmt"
    "errors"
    "encoding/json"    
    zmq "github.com/alecthomas/gozmq" 

)
type rpcMan struct { 
    Context *zmq.Context
    Socket *zmq.Socket
    ServAddr string 
} 
    

type Request struct {
    Method string
    Args interface{}
} 

type Response struct {
    Response interface{}
    Status int
} 

func Init(addr string) rpcMan { 
    context, _ := zmq.NewContext()
    socket, _ := context.NewSocket(zmq.REQ)

    return rpcMan{context,socket,addr} 
} 

func (rpc rpcMan) Connect() {
    fmt.Println("Connecting to rpc server at ",rpc.ServAddr) 
    rpc.Socket.Connect(rpc.ServAddr)
}     

func (rpc rpcMan) Close() {
    defer rpc.Context.Close()
    defer rpc.Socket.Close()
}  

func (rpc rpcMan) Call(method string, args ...interface{}) (interface{}, error){ 
    

    msg := Request{method, args} 

    enc, err := json.Marshal(msg)
    if err != nil {
        fmt.Println("json error", err) 
        return -1, err        
    } 
    rpc.Socket.Send(enc, 0)
    resp := new(Response) 
    
    // Wait for reply, could probably start a timeout thing
    reply, _ := rpc.Socket.Recv(0)
        
    err = json.Unmarshal(reply,resp) 
    
    if err != nil {
        fmt.Println("json error", err)
        return -1, err
    }

    if resp.Status != 0 {
        return -1, errors.New("RPC Method not found:"+method) 
    } 
    
    return resp.Response, nil

}

