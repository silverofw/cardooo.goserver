using System;
namespace cardooo.goserver
{
    public class ApiEndpointBasic : ApiEndpoint
    {
        public override void Excute(string param, Action<string> error = null)
        {
            base.Excute(param, error);
        }
    }
}
