using System;
using System.Collections;
using System.Collections.Generic;
using System.Net.Sockets;
using System.Text;
using UnityEngine;

public class GoServer : MonoBehaviour
{
    private TcpClient client;
    private NetworkStream stream;

    // 要連接的伺服器 IP 和 Port
    public string serverIp = "127.0.0.1";
    public int serverPort = 8000;

    // Start is called before the first frame update
    void Start()
    {
        try
        {
            // 建立一個 TCP 連線
            client = new TcpClient(serverIp, serverPort);
            stream = client.GetStream();
        }
        catch (Exception e)
        {
            Debug.Log("Error: " + e.Message);
        }
    }

    // Update is called once per frame
    void Update()
    {
        if (client == null)
        {
            return;
        }

        if (stream.DataAvailable)
        {
            // 讀取從伺服器接收到的消息
            byte[] data = new byte[1024];
            int bytes = stream.Read(data, 0, data.Length);
            string message = Encoding.ASCII.GetString(data, 0, bytes);
            Debug.Log("Received: " + message);
        }
    }

    [ContextMenu("sendToServer")]
    void sendToServer()
    {
        // 發送消息到伺服器
        string message = "Hello, server!";
        byte[] data = Encoding.ASCII.GetBytes(message);
        stream.Write(data, 0, data.Length);
        Debug.Log("Sent: " + message);
    }

    private void OnDestroy()
    {
        // 關閉 TCP 連線
        if (client != null)
        {
            client.Close();
        }
    }
}
