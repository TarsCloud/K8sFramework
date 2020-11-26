﻿/**
 * Tencent is pleased to support the open source community by making Tars available.
 *
 * Copyright (C) 2016THL A29 Limited, a Tencent company. All rights reserved.
 *
 * Licensed under the BSD 3-Clause License (the "License"); you may not use this file except 
 * in compliance with the License. You may obtain a copy of the License at
 *
 * https://opensource.org/licenses/BSD-3-Clause
 *
 * Unless required by applicable law or agreed to in writing, software distributed 
 * under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR 
 * CONDITIONS OF ANY KIND, either express or implied. See the License for the 
 * specific language governing permissions and limitations under the License.
 */
 
﻿// **********************************************************************
// This file was generated by a TARS parser!
// TARS version 1.1.0.
// **********************************************************************

"use strict";

var assert    = require("assert");
var TarsStream = require("@tars/stream");

var tars = tars || {};
module.exports.tars = tars;

tars.ServerState = {
    "Inactive" : 0,
    "Activating" : 1,
    "Active" : 2,
    "Deactivating" : 3,
    "Destroying" : 4,
    "Destroyed" : 5,
    "_classname" : "tars.ServerState"
};
tars.ServerState._write = function(os, tag, val) { return os.writeInt32(tag, val); };
tars.ServerState._read  = function(is, tag, def) { return is.readInt32(tag, true, def); };

tars.LoadInfo = function() {
    this.avg1 = 0;
    this.avg5 = 0;
    this.avg15 = 0;
    this.avgCpu = 0;
    this._classname = "tars.LoadInfo";
};
tars.LoadInfo._classname = "tars.LoadInfo";
tars.LoadInfo._write = function (os, tag, value) { os.writeStruct(tag, value); };
tars.LoadInfo._read  = function (is, tag, def) { return is.readStruct(tag, true, def); };
tars.LoadInfo._readFrom = function (is) {
    var tmp = new tars.LoadInfo();
    tmp.avg1 = is.readFloat(0, true, 0);
    tmp.avg5 = is.readFloat(1, true, 0);
    tmp.avg15 = is.readFloat(2, true, 0);
    tmp.avgCpu = is.readInt32(3, false, 0);
    return tmp;
};
tars.LoadInfo.prototype._writeTo = function (os) {
    os.writeFloat(0, this.avg1);
    os.writeFloat(1, this.avg5);
    os.writeFloat(2, this.avg15);
    os.writeInt32(3, this.avgCpu);
};
tars.LoadInfo.prototype._equal = function () {
    assert.fail("this structure not define key operation");
};
tars.LoadInfo.prototype._genKey = function () {
    if (!this._proto_struct_name_) {
        this._proto_struct_name_ = "STRUCT" + Math.random();
    }
    return this._proto_struct_name_;
};
tars.LoadInfo.prototype.toObject = function() { 
    return {
        "avg1" : this.avg1,
        "avg5" : this.avg5,
        "avg15" : this.avg15,
        "avgCpu" : this.avgCpu
    };
};
tars.LoadInfo.prototype.readFromObject = function(json) { 
    json.hasOwnProperty("avg1") && (this.avg1 = json.avg1);
    json.hasOwnProperty("avg5") && (this.avg5 = json.avg5);
    json.hasOwnProperty("avg15") && (this.avg15 = json.avg15);
    json.hasOwnProperty("avgCpu") && (this.avgCpu = json.avgCpu);
};
tars.LoadInfo.prototype.toBinBuffer = function () {
    var os = new TarsStream.TarsOutputStream();
    this._writeTo(os);
    return os.getBinBuffer();
};
tars.LoadInfo.new = function () {
    return new tars.LoadInfo();
};
tars.LoadInfo.create = function (is) {
    return tars.LoadInfo._readFrom(is);
};

tars.PatchInfo = function() {
    this.bPatching = false;
    this.iPercent = 0;
    this.iModifyTime = 0;
    this.sVersion = "";
    this.sResult = "";
    this.bSucc = false;
    this._classname = "tars.PatchInfo";
};
tars.PatchInfo._classname = "tars.PatchInfo";
tars.PatchInfo._write = function (os, tag, value) { os.writeStruct(tag, value); };
tars.PatchInfo._read  = function (is, tag, def) { return is.readStruct(tag, true, def); };
tars.PatchInfo._readFrom = function (is) {
    var tmp = new tars.PatchInfo();
    tmp.bPatching = is.readBoolean(0, true, false);
    tmp.iPercent = is.readInt32(1, true, 0);
    tmp.iModifyTime = is.readInt32(2, true, 0);
    tmp.sVersion = is.readString(3, true, "");
    tmp.sResult = is.readString(4, true, "");
    tmp.bSucc = is.readBoolean(5, false, false);
    return tmp;
};
tars.PatchInfo.prototype._writeTo = function (os) {
    os.writeBoolean(0, this.bPatching);
    os.writeInt32(1, this.iPercent);
    os.writeInt32(2, this.iModifyTime);
    os.writeString(3, this.sVersion);
    os.writeString(4, this.sResult);
    os.writeBoolean(5, this.bSucc);
};
tars.PatchInfo.prototype._equal = function () {
    assert.fail("this structure not define key operation");
};
tars.PatchInfo.prototype._genKey = function () {
    if (!this._proto_struct_name_) {
        this._proto_struct_name_ = "STRUCT" + Math.random();
    }
    return this._proto_struct_name_;
};
tars.PatchInfo.prototype.toObject = function() { 
    return {
        "bPatching" : this.bPatching,
        "iPercent" : this.iPercent,
        "iModifyTime" : this.iModifyTime,
        "sVersion" : this.sVersion,
        "sResult" : this.sResult,
        "bSucc" : this.bSucc
    };
};
tars.PatchInfo.prototype.readFromObject = function(json) { 
    json.hasOwnProperty("bPatching") && (this.bPatching = json.bPatching);
    json.hasOwnProperty("iPercent") && (this.iPercent = json.iPercent);
    json.hasOwnProperty("iModifyTime") && (this.iModifyTime = json.iModifyTime);
    json.hasOwnProperty("sVersion") && (this.sVersion = json.sVersion);
    json.hasOwnProperty("sResult") && (this.sResult = json.sResult);
    json.hasOwnProperty("bSucc") && (this.bSucc = json.bSucc);
};
tars.PatchInfo.prototype.toBinBuffer = function () {
    var os = new TarsStream.TarsOutputStream();
    this._writeTo(os);
    return os.getBinBuffer();
};
tars.PatchInfo.new = function () {
    return new tars.PatchInfo();
};
tars.PatchInfo.create = function (is) {
    return tars.PatchInfo._readFrom(is);
};

tars.PreparePatchInfo = function() {
    this.bPreparePatching = false;
    this.iPercent = 0;
    this.iModifyTime = 0;
    this.sVersion = "";
    this.sResult = "";
    this.ret = 0;
    this._classname = "tars.PreparePatchInfo";
};
tars.PreparePatchInfo._classname = "tars.PreparePatchInfo";
tars.PreparePatchInfo._write = function (os, tag, value) { os.writeStruct(tag, value); };
tars.PreparePatchInfo._read  = function (is, tag, def) { return is.readStruct(tag, true, def); };
tars.PreparePatchInfo._readFrom = function (is) {
    var tmp = new tars.PreparePatchInfo();
    tmp.bPreparePatching = is.readBoolean(0, true, false);
    tmp.iPercent = is.readInt32(1, true, 0);
    tmp.iModifyTime = is.readInt32(2, true, 0);
    tmp.sVersion = is.readString(3, true, "");
    tmp.sResult = is.readString(4, true, "");
    tmp.ret = is.readInt32(5, false, 0);
    return tmp;
};
tars.PreparePatchInfo.prototype._writeTo = function (os) {
    os.writeBoolean(0, this.bPreparePatching);
    os.writeInt32(1, this.iPercent);
    os.writeInt32(2, this.iModifyTime);
    os.writeString(3, this.sVersion);
    os.writeString(4, this.sResult);
    os.writeInt32(5, this.ret);
};
tars.PreparePatchInfo.prototype._equal = function () {
    assert.fail("this structure not define key operation");
};
tars.PreparePatchInfo.prototype._genKey = function () {
    if (!this._proto_struct_name_) {
        this._proto_struct_name_ = "STRUCT" + Math.random();
    }
    return this._proto_struct_name_;
};
tars.PreparePatchInfo.prototype.toObject = function() { 
    return {
        "bPreparePatching" : this.bPreparePatching,
        "iPercent" : this.iPercent,
        "iModifyTime" : this.iModifyTime,
        "sVersion" : this.sVersion,
        "sResult" : this.sResult,
        "ret" : this.ret
    };
};
tars.PreparePatchInfo.prototype.readFromObject = function(json) { 
    json.hasOwnProperty("bPreparePatching") && (this.bPreparePatching = json.bPreparePatching);
    json.hasOwnProperty("iPercent") && (this.iPercent = json.iPercent);
    json.hasOwnProperty("iModifyTime") && (this.iModifyTime = json.iModifyTime);
    json.hasOwnProperty("sVersion") && (this.sVersion = json.sVersion);
    json.hasOwnProperty("sResult") && (this.sResult = json.sResult);
    json.hasOwnProperty("ret") && (this.ret = json.ret);
};
tars.PreparePatchInfo.prototype.toBinBuffer = function () {
    var os = new TarsStream.TarsOutputStream();
    this._writeTo(os);
    return os.getBinBuffer();
};
tars.PreparePatchInfo.new = function () {
    return new tars.PreparePatchInfo();
};
tars.PreparePatchInfo.create = function (is) {
    return tars.PreparePatchInfo._readFrom(is);
};

tars.NodeInfo = function() {
    this.nodeName = "";
    this.nodeObj = "";
    this.endpointIp = "";
    this.endpointPort = 0;
    this.timeOut = 0;
    this.dataDir = "";
    this.version = "";
    this.coreFileSize = "";
    this.openFiles = 0;
    this._classname = "tars.NodeInfo";
};
tars.NodeInfo._classname = "tars.NodeInfo";
tars.NodeInfo._write = function (os, tag, value) { os.writeStruct(tag, value); };
tars.NodeInfo._read  = function (is, tag, def) { return is.readStruct(tag, true, def); };
tars.NodeInfo._readFrom = function (is) {
    var tmp = new tars.NodeInfo();
    tmp.nodeName = is.readString(0, true, "");
    tmp.nodeObj = is.readString(1, true, "");
    tmp.endpointIp = is.readString(2, true, "");
    tmp.endpointPort = is.readInt32(3, true, 0);
    tmp.timeOut = is.readInt16(4, true, 0);
    tmp.dataDir = is.readString(5, true, "");
    tmp.version = is.readString(6, false, "");
    tmp.coreFileSize = is.readString(7, false, "");
    tmp.openFiles = is.readInt32(8, false, 0);
    return tmp;
};
tars.NodeInfo.prototype._writeTo = function (os) {
    os.writeString(0, this.nodeName);
    os.writeString(1, this.nodeObj);
    os.writeString(2, this.endpointIp);
    os.writeInt32(3, this.endpointPort);
    os.writeInt16(4, this.timeOut);
    os.writeString(5, this.dataDir);
    os.writeString(6, this.version);
    os.writeString(7, this.coreFileSize);
    os.writeInt32(8, this.openFiles);
};
tars.NodeInfo.prototype._equal = function () {
    assert.fail("this structure not define key operation");
};
tars.NodeInfo.prototype._genKey = function () {
    if (!this._proto_struct_name_) {
        this._proto_struct_name_ = "STRUCT" + Math.random();
    }
    return this._proto_struct_name_;
};
tars.NodeInfo.prototype.toObject = function() { 
    return {
        "nodeName" : this.nodeName,
        "nodeObj" : this.nodeObj,
        "endpointIp" : this.endpointIp,
        "endpointPort" : this.endpointPort,
        "timeOut" : this.timeOut,
        "dataDir" : this.dataDir,
        "version" : this.version,
        "coreFileSize" : this.coreFileSize,
        "openFiles" : this.openFiles
    };
};
tars.NodeInfo.prototype.readFromObject = function(json) { 
    json.hasOwnProperty("nodeName") && (this.nodeName = json.nodeName);
    json.hasOwnProperty("nodeObj") && (this.nodeObj = json.nodeObj);
    json.hasOwnProperty("endpointIp") && (this.endpointIp = json.endpointIp);
    json.hasOwnProperty("endpointPort") && (this.endpointPort = json.endpointPort);
    json.hasOwnProperty("timeOut") && (this.timeOut = json.timeOut);
    json.hasOwnProperty("dataDir") && (this.dataDir = json.dataDir);
    json.hasOwnProperty("version") && (this.version = json.version);
    json.hasOwnProperty("coreFileSize") && (this.coreFileSize = json.coreFileSize);
    json.hasOwnProperty("openFiles") && (this.openFiles = json.openFiles);
};
tars.NodeInfo.prototype.toBinBuffer = function () {
    var os = new TarsStream.TarsOutputStream();
    this._writeTo(os);
    return os.getBinBuffer();
};
tars.NodeInfo.new = function () {
    return new tars.NodeInfo();
};
tars.NodeInfo.create = function (is) {
    return tars.NodeInfo._readFrom(is);
};

tars.ServerStateInfo = function() {
    this.serverState = tars.ServerState.Inactive;
    this.processId = 0;
    this.nodeName = "";
    this.application = "";
    this.serverName = "";
    this.settingState = tars.ServerState.Inactive;
    this._classname = "tars.ServerStateInfo";
};
tars.ServerStateInfo._classname = "tars.ServerStateInfo";
tars.ServerStateInfo._write = function (os, tag, value) { os.writeStruct(tag, value); };
tars.ServerStateInfo._read  = function (is, tag, def) { return is.readStruct(tag, true, def); };
tars.ServerStateInfo._readFrom = function (is) {
    var tmp = new tars.ServerStateInfo();
    tmp.serverState = is.readInt32(0, true, tars.ServerState.Inactive);
    tmp.processId = is.readInt32(1, true, 0);
    tmp.nodeName = is.readString(2, false, "");
    tmp.application = is.readString(3, false, "");
    tmp.serverName = is.readString(4, false, "");
    tmp.settingState = is.readInt32(5, false, tars.ServerState.Inactive);
    return tmp;
};
tars.ServerStateInfo.prototype._writeTo = function (os) {
    os.writeInt32(0, this.serverState);
    os.writeInt32(1, this.processId);
    os.writeString(2, this.nodeName);
    os.writeString(3, this.application);
    os.writeString(4, this.serverName);
    os.writeInt32(5, this.settingState);
};
tars.ServerStateInfo.prototype._equal = function () {
    assert.fail("this structure not define key operation");
};
tars.ServerStateInfo.prototype._genKey = function () {
    if (!this._proto_struct_name_) {
        this._proto_struct_name_ = "STRUCT" + Math.random();
    }
    return this._proto_struct_name_;
};
tars.ServerStateInfo.prototype.toObject = function() { 
    return {
        "serverState" : this.serverState,
        "processId" : this.processId,
        "nodeName" : this.nodeName,
        "application" : this.application,
        "serverName" : this.serverName,
        "settingState" : this.settingState
    };
};
tars.ServerStateInfo.prototype.readFromObject = function(json) { 
    json.hasOwnProperty("serverState") && (this.serverState = json.serverState);
    json.hasOwnProperty("processId") && (this.processId = json.processId);
    json.hasOwnProperty("nodeName") && (this.nodeName = json.nodeName);
    json.hasOwnProperty("application") && (this.application = json.application);
    json.hasOwnProperty("serverName") && (this.serverName = json.serverName);
    json.hasOwnProperty("settingState") && (this.settingState = json.settingState);
};
tars.ServerStateInfo.prototype.toBinBuffer = function () {
    var os = new TarsStream.TarsOutputStream();
    this._writeTo(os);
    return os.getBinBuffer();
};
tars.ServerStateInfo.new = function () {
    return new tars.ServerStateInfo();
};
tars.ServerStateInfo.create = function (is) {
    return tars.ServerStateInfo._readFrom(is);
};

tars.PatchRequest = function() {
    this.appname = "";
    this.servername = "";
    this.nodename = "";
    this.groupname = "";
    this.binname = "";
    this.version = "";
    this.user = "";
    this.servertype = "";
    this.patchobj = "";
    this.md5 = "";
    this.ostype = "";
    this.filepath = "";
    this._classname = "tars.PatchRequest";
};
tars.PatchRequest._classname = "tars.PatchRequest";
tars.PatchRequest._write = function (os, tag, value) { os.writeStruct(tag, value); };
tars.PatchRequest._read  = function (is, tag, def) { return is.readStruct(tag, true, def); };
tars.PatchRequest._readFrom = function (is) {
    var tmp = new tars.PatchRequest();
    tmp.appname = is.readString(0, true, "");
    tmp.servername = is.readString(1, true, "");
    tmp.nodename = is.readString(2, true, "");
    tmp.groupname = is.readString(3, true, "");
    tmp.binname = is.readString(4, true, "");
    tmp.version = is.readString(5, true, "");
    tmp.user = is.readString(6, true, "");
    tmp.servertype = is.readString(7, true, "");
    tmp.patchobj = is.readString(8, true, "");
    tmp.md5 = is.readString(9, true, "");
    tmp.ostype = is.readString(10, false, "");
    tmp.filepath = is.readString(11, false, "");
    return tmp;
};
tars.PatchRequest.prototype._writeTo = function (os) {
    os.writeString(0, this.appname);
    os.writeString(1, this.servername);
    os.writeString(2, this.nodename);
    os.writeString(3, this.groupname);
    os.writeString(4, this.binname);
    os.writeString(5, this.version);
    os.writeString(6, this.user);
    os.writeString(7, this.servertype);
    os.writeString(8, this.patchobj);
    os.writeString(9, this.md5);
    os.writeString(10, this.ostype);
    os.writeString(11, this.filepath);
};
tars.PatchRequest.prototype._equal = function () {
    assert.fail("this structure not define key operation");
};
tars.PatchRequest.prototype._genKey = function () {
    if (!this._proto_struct_name_) {
        this._proto_struct_name_ = "STRUCT" + Math.random();
    }
    return this._proto_struct_name_;
};
tars.PatchRequest.prototype.toObject = function() { 
    return {
        "appname" : this.appname,
        "servername" : this.servername,
        "nodename" : this.nodename,
        "groupname" : this.groupname,
        "binname" : this.binname,
        "version" : this.version,
        "user" : this.user,
        "servertype" : this.servertype,
        "patchobj" : this.patchobj,
        "md5" : this.md5,
        "ostype" : this.ostype,
        "filepath" : this.filepath
    };
};
tars.PatchRequest.prototype.readFromObject = function(json) { 
    json.hasOwnProperty("appname") && (this.appname = json.appname);
    json.hasOwnProperty("servername") && (this.servername = json.servername);
    json.hasOwnProperty("nodename") && (this.nodename = json.nodename);
    json.hasOwnProperty("groupname") && (this.groupname = json.groupname);
    json.hasOwnProperty("binname") && (this.binname = json.binname);
    json.hasOwnProperty("version") && (this.version = json.version);
    json.hasOwnProperty("user") && (this.user = json.user);
    json.hasOwnProperty("servertype") && (this.servertype = json.servertype);
    json.hasOwnProperty("patchobj") && (this.patchobj = json.patchobj);
    json.hasOwnProperty("md5") && (this.md5 = json.md5);
    json.hasOwnProperty("ostype") && (this.ostype = json.ostype);
    json.hasOwnProperty("filepath") && (this.filepath = json.filepath);
};
tars.PatchRequest.prototype.toBinBuffer = function () {
    var os = new TarsStream.TarsOutputStream();
    this._writeTo(os);
    return os.getBinBuffer();
};
tars.PatchRequest.new = function () {
    return new tars.PatchRequest();
};
tars.PatchRequest.create = function (is) {
    return tars.PatchRequest._readFrom(is);
};

tars.PreparePatchRequest = function() {
    this.appname = "";
    this.servername = "";
    this.groupname = "";
    this.version = "";
    this.user = "";
    this.servertype = "";
    this.patchobj = "";
    this.md5 = "";
    this.ostype = "";
    this.specialNodeList = new TarsStream.List(TarsStream.String);
    this.filepath = "";
    this._classname = "tars.PreparePatchRequest";
};
tars.PreparePatchRequest._classname = "tars.PreparePatchRequest";
tars.PreparePatchRequest._write = function (os, tag, value) { os.writeStruct(tag, value); };
tars.PreparePatchRequest._read  = function (is, tag, def) { return is.readStruct(tag, true, def); };
tars.PreparePatchRequest._readFrom = function (is) {
    var tmp = new tars.PreparePatchRequest();
    tmp.appname = is.readString(0, true, "");
    tmp.servername = is.readString(1, true, "");
    tmp.groupname = is.readString(2, true, "");
    tmp.version = is.readString(3, true, "");
    tmp.user = is.readString(4, true, "");
    tmp.servertype = is.readString(5, true, "");
    tmp.patchobj = is.readString(6, true, "");
    tmp.md5 = is.readString(7, true, "");
    tmp.ostype = is.readString(8, true, "");
    tmp.specialNodeList = is.readList(9, true, TarsStream.List(TarsStream.String));
    tmp.filepath = is.readString(10, false, "");
    return tmp;
};
tars.PreparePatchRequest.prototype._writeTo = function (os) {
    os.writeString(0, this.appname);
    os.writeString(1, this.servername);
    os.writeString(2, this.groupname);
    os.writeString(3, this.version);
    os.writeString(4, this.user);
    os.writeString(5, this.servertype);
    os.writeString(6, this.patchobj);
    os.writeString(7, this.md5);
    os.writeString(8, this.ostype);
    os.writeList(9, this.specialNodeList);
    os.writeString(10, this.filepath);
};
tars.PreparePatchRequest.prototype._equal = function (anItem) {
    return this.appname === anItem.appname && 
        this.servername === anItem.servername && 
        this.version === anItem.version;
};
tars.PreparePatchRequest.prototype._genKey = function () {
    if (!this._proto_struct_name_) {
        this._proto_struct_name_ = "STRUCT" + Math.random();
    }
    return this._proto_struct_name_;
};
tars.PreparePatchRequest.prototype.toObject = function() { 
    return {
        "appname" : this.appname,
        "servername" : this.servername,
        "groupname" : this.groupname,
        "version" : this.version,
        "user" : this.user,
        "servertype" : this.servertype,
        "patchobj" : this.patchobj,
        "md5" : this.md5,
        "ostype" : this.ostype,
        "specialNodeList" : this.specialNodeList.toObject(),
        "filepath" : this.filepath
    };
};
tars.PreparePatchRequest.prototype.readFromObject = function(json) { 
    json.hasOwnProperty("appname") && (this.appname = json.appname);
    json.hasOwnProperty("servername") && (this.servername = json.servername);
    json.hasOwnProperty("groupname") && (this.groupname = json.groupname);
    json.hasOwnProperty("version") && (this.version = json.version);
    json.hasOwnProperty("user") && (this.user = json.user);
    json.hasOwnProperty("servertype") && (this.servertype = json.servertype);
    json.hasOwnProperty("patchobj") && (this.patchobj = json.patchobj);
    json.hasOwnProperty("md5") && (this.md5 = json.md5);
    json.hasOwnProperty("ostype") && (this.ostype = json.ostype);
    json.hasOwnProperty("specialNodeList") && (this.specialNodeList.readFromObject(json.specialNodeList));
    json.hasOwnProperty("filepath") && (this.filepath = json.filepath);
};
tars.PreparePatchRequest.prototype.toBinBuffer = function () {
    var os = new TarsStream.TarsOutputStream();
    this._writeTo(os);
    return os.getBinBuffer();
};
tars.PreparePatchRequest.new = function () {
    return new tars.PreparePatchRequest();
};
tars.PreparePatchRequest.create = function (is) {
    return tars.PreparePatchRequest._readFrom(is);
};



