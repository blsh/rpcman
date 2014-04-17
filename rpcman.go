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
    return RPCMan{context,addr} 
} 


func (rpc RPCMan) Close() {
    defer rpc.Context.Close()
}  

func (rpc RPCMan) Call(method string, args ...interface{}) (interface{}, error){ 
    

    msg := Request{method, args} 

    enc, err := json.Marshal(msg)
    if err != nil {
        log.Println("json error", err) 
        return -1, err        
    } 
    
    socket, err := rpc.Context.NewSocket(zmq.REQ)
    if err != nil {
        log.Println("could not set up socket to rpcman",err)
        return -1, err
    }

    // set timeout to 1s
    duration, _ := time.ParseDuration("1s")
    socket.SetRcvTimeout(duration)
    
    err = socket.Connect(rpc.ServAddr)
    if err != nil {
        log.Println("Could not connect to socket",err)
        return -1, err
    }
    

    
    socket.Send(enc,0)
    //rpc.Socket.Send(enc, 0)
    resp := new(Response) 
    
    //reply, err := rpc.Socket.Recv(0)
    reply, err := socket.Recv(0)
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

    socket.Close()
    

    
    return resp.Response, nil

}

