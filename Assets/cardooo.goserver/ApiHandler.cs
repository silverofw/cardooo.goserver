using System;
using System.Collections.Generic;
namespace cardooo.goserver
{
    public class ApiHandler
    {
        static ApiHandler inst;
        public static ApiHandler Inst => inst ??= new ApiHandler();

        Queue<ResponsePack> responseQueue = new Queue<ResponsePack>();
        Dictionary<int, ApiEndpoint> apiDic = new Dictionary<int, ApiEndpoint>();
        public void addApi(int apiId, ApiEndpoint apiEndpoint)
        {
            if (apiDic.ContainsKey(apiId))
                return;
            apiDic.Add(apiId, apiEndpoint);
        }

        public void AddResponsePack(int systemId, int apiId, string param, Action<string> error = null)
        {
            responseQueue.Enqueue(new ResponsePack() { systemId = systemId, apiId = apiId, param = param });
        }

        public void ProcessRespone(Action<string> error = null)
        {
            if (responseQueue.Count == 0)
            {
                return;
            }
            var pack = responseQueue.Dequeue();
            if (!apiDic.TryGetValue(pack.apiId, out var apiEndpoint))
            {
                error?.Invoke($"[{pack.apiId}] can not find api~");
                return;
            }
            apiEndpoint.Excute(pack.param, error);
        }
    }
}
