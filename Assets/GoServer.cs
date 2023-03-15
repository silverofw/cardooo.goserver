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

    void Start()
    {
        if (!nonThread)
        {
            GoServerHandler.Inst.JoinServer(serverIp, serverPort);
        }
        else
        {
            GoServerHandler.Inst.JoinServerNonThread(serverIp, serverPort);
        }
        ApiHandler.Inst.addApi(1, new ApiEndpointLog());
        ApiHandler.Inst.addApi(2, new ApiEndpointLog());
        ApiHandler.Inst.addApi(3, new ApiEndpointBoardcast());
    }

    private void Update()
    {
        GoServerHandler.Inst.updateStream((error) => {
            Debug.Log(error);
        });
    }

    [ContextMenu("sendMsg")]
    void sendMsg()
    {
        var msg = $"{uid:0000}{systemIndex:0000}{apiIndex:0000}{param}";
        GoServerHandler.Inst.sendToServer(msg);
        Debug.Log("Sent: " + msg);
    }

    private void OnDestroy()
    {
        GoServerHandler.Inst.LeaveServer();
    }
}
