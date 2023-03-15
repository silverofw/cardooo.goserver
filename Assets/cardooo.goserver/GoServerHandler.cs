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
        /// <summary>
        /// use thread update or not
        /// </summary>
        public bool NonThread { get; private set; } = false;

        private TcpClient client;
        private NetworkStream stream;
        private Thread thread;

        public void JoinServer(string ip, int port, Action<string> error = null)
        {
            if (client != null && client.Connected)
            {
                return;
            }

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

                string systemId = Encoding.ASCII.GetString(data, 0, 4);
                string apiId = Encoding.ASCII.GetString(data, 4, 4);
                string param = Encoding.ASCII.GetString(data, 8, bytes - 8);

                ApiHandler.Inst.AddResponsePack(int.Parse(systemId), int.Parse(apiId), param);
            }
        }

        public void sendToServer(string message)
        {
            byte[] data = Encoding.ASCII.GetBytes(message);
            stream.Write(data, 0, data.Length);
        }

        public void LeaveServer()
        {
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
