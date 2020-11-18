const Taf = require("@taf/taf-rpc").client;
const ldapPrx = require("../proxy/LDAPProxy").Common;

let adaptorStr = 'Comm.LDAPServer.LDAPObj'
if (!process.env.TAF_CONFIG) {
    adaptorStr += '@tcp -h 172.16.8.125 -t 60000 -p 8888'
}
const prxObj = Taf.stringToProxy(ldapPrx.LDAPServerProxy, adaptorStr);

exports = module.exports = {
    // 校验ldapkey合法性
    checkLDAPKey(ldapKey) {
        return new Promise(async (resolve, reject) => {
            const req = new ldapPrx.CheckLDAPKeyReq()
            req.readFromObject({ LDAPKey: ldapKey })
            try{
                const result = await prxObj.checkLDAPKey(req)
                resolve(result.response.arguments.rsp)
            }catch(err) {
                reject(new Error(err))
            }
        })
    },
    // 退出登录
    logoutLDAPKey(ldapKey) {
        return new Promise(async (resolve, reject) => {
            const req = new ldapPrx.LogoutReq()
            req.ldapKey = ldapKey
            try{
                const result = await prxObj.logout(req)
                resolve(result.response.arguments.rsp)
            }catch(err) {
                reject(new Error(err))
            }
        })
    }
}