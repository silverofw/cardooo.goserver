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

    // �n�s�������A�� IP �M Port
    public string serverIp = "127.0.0.1";
    public int serverPort = 8000;

    // Start is called before the first frame update
    void Start()
    {
        try
        {
            // �إߤ@�� TCP �s�u
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
            // Ū���q���A�������쪺����
            byte[] data = new byte[1024];
            int bytes = stream.Read(data, 0, data.Length);
            string message = Encoding.ASCII.GetString(data, 0, bytes);
            Debug.Log("Received: " + message);
        }
    }

    [ContextMenu("sendToServer")]
    void sendToServer()
    {
        // �o�e��������A��
        string message = "Hello, server!";
        byte[] data = Encoding.ASCII.GetBytes(message);
        stream.Write(data, 0, data.Length);
        Debug.Log("Sent: " + message);
    }

    private void OnDestroy()
    {
        // ���� TCP �s�u
        if (client != null)
        {
            client.Close();
        }
    }
}
