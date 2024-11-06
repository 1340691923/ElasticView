export function SaveEsConnect(link){
  sessionStorage.setItem('EsConnect', link)
}

export function GetEsConnect(){
  let esConnect = sessionStorage.getItem('EsConnect')
  if(esConnect) {
    return Number(esConnect)
  }
  if(esConnect == null){
    esConnect = 0
  }
  return esConnect
}

export function SaveEsConnectVer(ver){
  sessionStorage.setItem('EsConnVersion', ver)
}

export function GetEsConnectVer(){
  let esConnectVer = sessionStorage.getItem('EsConnVersion')
  if(esConnectVer) {
    return esConnectVer
  }
  if(esConnectVer == null || esConnectVer == undefined){
    esConnectVer = ""
  }
  return esConnectVer
}
