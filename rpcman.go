package rpcman

import (
    "log"
    "errors"
    "encoding/json"    
    "time"
    zmq "github.com/alecthomas/gozmq" 

)
type RPCMan struct { 
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

func Init(addr string) RPCMan { 
    context, _ := zmq.NewContext()
    socket, _ := context.NewSocket(zmq.REQ)

    // set timeout to 1s
    duration, _ := time.ParseDuration("1s")
    socket.SetRcvTimeout(duration)

    return RPCMan{context,socket,addr} 
} 

func (rpc RPCMan) Connect() error {
    log.Println("Connecting to rpc server at ",rpc.ServAddr) 
    return rpc.Socket.Connect(rpc.ServAddr)
}     

func (rpc RPCMan) Close() {
    defer rpc.Context.Close()
    defer rpc.Socket.Close()
}  

func (rpc RPCMan) Call(method string, args ...interface{}) (interface{}, error){ 
    

    msg := Request{method, args} 

    enc, err := json.Marshal(msg)
    if err != nil {
        log.Println("json error", err) 
        return -1, err        
    } 
    rpc.Socket.Send(enc, 0)
    resp := new(Response) 
    
    reply, err := rpc.Socket.Recv(0)
    if err != nil {
        log.Println("RPC recv failed:", err)
        return -1, err
    }
        
    err = json.Unmarshal(reply,resp) 
    
    if err != nil {
        log.Println("json error", err)
        return -1, err
    }

    if resp.Status != 0 {
        return -1, errors.New("RPC Method not found:"+method) 
    } 
    
    return resp.Response, nil

}

