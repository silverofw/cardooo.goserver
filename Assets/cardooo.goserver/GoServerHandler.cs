using System;
using System.Net.Sockets;
using System.Text;
using System.Threading;

namespace cardooo.goserver
{
    public class GoServerHandler
    {
        static GoServerHandler inst;
        public static GoServerHandler Inst => inst ??= new GoServerHandler();
        public string serverIp { get; private set; } = "127.0.0.1";
        public int serverPort { get; private set; } = 8000;
        public int packSize { get; private set; } = 1024;
        public bool IsConnected { get {
                if (client != null && client.Connected)
                    return true;
                return false;
            }
        }

        /// <summary>
        /// use thread update or not
        /// </summary>
        public bool NonThread { get; private set; } = false;

        private TcpClient client;
        private NetworkStream stream;
        private Thread thread;

        string getStream = "";

        Action<string> error;

        public void JoinServer(string ip, int port, Action<string> error = null)
        {
            if (client != null && client.Connected)
            {
                return;
            }
            this.error = error;

            serverIp = ip;
            serverPort = port;
            NonThread = false;

            try
            {
                client = new TcpClient(serverIp, serverPort);
                stream = client.GetStream();

                if (thread != null)
                {
                    thread.Abort();
                }
                thread = new Thread(new ThreadStart(updateStreamThread));
                thread.Start();
            }
            catch (Exception e)
            {
                error?.Invoke("Error: " + e.Message);
            }
        }


        public void JoinServerNonThread(string ip, int port, Action<string> error = null)
        {
            if (client != null && client.Connected)
            {
                return;
            }
            this.error = error;

            serverIp = ip;
            serverPort = port;
            NonThread = true;

            try
            {
                client = new TcpClient(serverIp, serverPort);
                stream = client.GetStream();
            }
            catch (Exception e)
            {
                error?.Invoke("Error: " + e.Message);
            }
        }

        void updateStreamThread()
        {
            while (client != null && client.Connected)
            {
                onSteamDataAvailable();
            }
            error("Server ShutDown!");
        }

        public void updateStream(Action<string> error = null)
        {
            ApiHandler.Inst.ProcessRespone(error);

            // only nonThread need call
            if (NonThread)
            {
                onSteamDataAvailable();
            }
        }

        void onSteamDataAvailable()
        {
            if (stream.DataAvailable)
            {
                byte[] data = new byte[packSize];
                int bytes = stream.Read(data, 0, data.Length);
                getStream += Encoding.ASCII.GetString(data);

                string[] strs = getStream.Split("[<]");
                int index = 1;
                while (strs.Length > index)
                {
                    bool isFinish = strs[index].Contains("[>]");
                    if (isFinish)
                    {
                        // 完整封包
                        string systemId = strs[index].Substring(0, 4);
                        string apiId = strs[index].Substring(4, 4);
                        string param = strs[index].Substring(8, strs[index].Length - 8);
                        param = param.Split("[>]")[0];
                        ApiHandler.Inst.AddResponsePack(int.Parse(systemId), int.Parse(apiId), param);
                        index++;
                    }
                    else
                    {
                        //封包不完整,繼續監聽
                        getStream = strs[index];
                    }
                }

                getStream = "";
            }
        }

        public void sendToServer(string message)
        {
            byte[] data = Encoding.ASCII.GetBytes(message);
            stream.Write(data, 0, data.Length);
        }

        public void ServerShutDown()
        {
            if (error != null)
                error("Server ShutDown!");
            LeaveServer();
        }

        public void LeaveServer()
        {
            error = null;
            if (thread != null)
            {
                thread.Abort();
            }
            if (client != null)
            {
                client.Close();
            }
        }
    }
}
