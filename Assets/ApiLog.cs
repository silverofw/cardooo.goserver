using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using cardooo.goserver;
using System;

public class ApiLog : ApiEndpoint
{
    public override void Excute(string param, Action<string> error = null)
    {
        base.Excute(param, error);

        Debug.Log($"[ApiEndpointLog] {param}");
    }
}
