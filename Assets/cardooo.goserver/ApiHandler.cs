using System;
using System.Collections.Generic;
namespace cardooo.goserver
{
    public class ApiHandler
    {
        static ApiHandler inst;
        public static ApiHandler Inst {
            get {
                if (inst == null)
                {
                    inst = new ApiHandler();
                }
                return inst;
            }
        }

        Dictionary<int, ApiEndpoint> apiDic = new Dictionary<int, ApiEndpoint>();

        public void addApi(int apiId, ApiEndpoint apiEndpoint)
        {
            if (apiDic.ContainsKey(apiId))
                return;
            apiDic.Add(apiId, apiEndpoint);
        }

        public void Response(int systemId, int apiId, string param, Action<string> error = null)
        {
            if (!apiDic.TryGetValue(apiId, out var apiEndpoint))
            {
                error?.Invoke($"[{apiId}] can not find api~");
                return;
            }
            apiEndpoint.Excute(param, error);
        }
    }
}
