using System.Collections.Generic;
using UnityEngine;
using cardooo.goserver;

public class GoServerTest : MonoBehaviour
{
    public string serverIp = "127.0.0.1";
    public int serverPort = 8000;
    public bool nonThread = false;

    public int uid;
    public int systemIndex;
    public int apiIndex;
    public string param;

    public List<string> errLogs = new List<string>();

    void Start()
    {
        ApiHandler.Inst.addApi(1, new ApiLog());
        ApiHandler.Inst.addApi(2, new ApiLog());
        ApiHandler.Inst.addApi(3, new ApiBoardcast());
    }

    void errorLog(string error)
    {
        errLogs.Add(error);
        Debug.LogError(error);
    }    

    private void Update()
    {
        GoServerHandler.Inst.updateStream((error) => {
            Debug.Log(error);
        });
    }

    [ContextMenu("connServer")]
    void connServer()
    {
        if (!nonThread)
        {
            GoServerHandler.Inst.JoinServer(serverIp, serverPort, errorLog);
        }
        else
        {
            GoServerHandler.Inst.JoinServerNonThread(serverIp, serverPort, errorLog);
        }
    }

    [ContextMenu("sendMsg")]
    void sendMsg()
    {
        if (!GoServerHandler.Inst.IsConnected)
        {
            Debug.LogError("Server is not Connected!");
            return;
        }
        var msg = $"{uid:0000}{systemIndex:0000}{apiIndex:0000}{param}";
        GoServerHandler.Inst.sendToServer(msg);
        Debug.Log("Sent: " + msg);
    }

    [ContextMenu("leaveServer")]
    void leaveServer() 
    { 
        GoServerHandler.Inst.LeaveServer();
    }

    private void OnDestroy()
    {
        leaveServer();
    }
}
