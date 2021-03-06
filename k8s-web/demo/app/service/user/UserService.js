let UserDao = require('../../dao/UserDao');
let logger = require('../../logger');
let sha1 = require('sha1');

const UserService = {}

UserService.getUserIdList = async() => {
    return await UserDao.getUserIdList();
}; 

//注册操作
UserService.modifyPass = async(uid, password) => {
    let updated = await UserDao.modifyPass(uid, sha1(password));
    return {updated: updated[0]};
};

UserService.addUser = async(userList) => {
    let newUserList = [];
    userList.forEach((user)=> {
        newUserList.push({
            uid: user.uid,
            password: sha1(user.password),
            update_time: new Date()
        });
    });
    return UserDao.insertUser(newUserList);
};

UserService.deleteUser = async(id) => {
    await UserDao.deleteUser({ 
        id
        });
};

module.exports = UserService;
