using System;
using System.Collections;
using System.Collections.Generic;
using System.Net.Sockets;
using System.Text;
using UnityEngine;
using cardooo.goserver;

public class GoServer : MonoBehaviour
{
    private TcpClient client;
    private NetworkStream stream;

    public string serverIp = "127.0.0.1";
    public int serverPort = 8000;

    public int uid;
    public int systemIndex;
    public int apiIndex;
    public string param;

    void Start()
    {
        try
        {
            client = new TcpClient(serverIp, serverPort);
            stream = client.GetStream();
        }
        catch (Exception e)
        {
            Debug.Log("Error: " + e.Message);
        }

        ApiHandler.Inst.addApi(1, new ApiEndpointLog());
        ApiHandler.Inst.addApi(2, new ApiEndpointLog());
        ApiHandler.Inst.addApi(3, new ApiEndpointBoardcast());
    }
    
    void Update()
    {
        if (client == null)
        {
            return;
        }

        if (stream.DataAvailable)
        {
            byte[] data = new byte[4096];
            int bytes = stream.Read(data, 0, data.Length);
            string message = Encoding.ASCII.GetString(data, 0, bytes);
            Debug.Log($"[{bytes}]Received: " + message);
            ApiHandler.Inst.Response(
                int.Parse(Encoding.ASCII.GetString(data, 0, 4)),
                int.Parse(Encoding.ASCII.GetString(data, 4, 4)),
                Encoding.ASCII.GetString(data, 8, bytes - 8), 
                (error)=>{
                    Debug.LogError(error);
            });
        }
    }

    [ContextMenu("setUid")]
    void setUid()
    {
        sendToServer($"{uid:0000}{systemIndex:0000}{1:0000}{uid:0000}");
    }

    [ContextMenu("sendMsg")]
    void sendMsg()
    {
        sendToServer($"{uid:0000}{systemIndex:0000}{apiIndex:0000}{param}");
    }

    void sendToServer(string message)
    {        
        byte[] data = Encoding.ASCII.GetBytes(message);
        stream.Write(data, 0, data.Length);
        Debug.Log("Sent: " + message);
    }

    private void OnDestroy()
    {
        if (client != null)
        {
            client.Close();
        }
    }
}
