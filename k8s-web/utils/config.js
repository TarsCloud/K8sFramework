/*
 * Tars config 读取
 */

'use strict';

const Q = require('q');
const ConfigParser = require('@tars/tars-utils').Config;

const config = {};

function parseConf(content, configFormat) {
    let ret = content;
    if (configFormat === 'c') {
        const configParser = new ConfigParser();
        configParser.parseText(content, 'utf8');
        ret = configParser.data;
    } else if (configFormat === 'json') {
        ret = JSON.parse(content);
    }
    return ret;
}

function loadConfig(filename, configFormat) {
    const dfd = Q.defer();
    if (process.env.TARS_CONFIG) {
        const tarsConfigHelper = require('@tars/tars-config');
        const helper = new tarsConfigHelper();
        helper.loadConfig(filename, configFormat).then((data) => {
            data = parseConf(data, configFormat);
            global.CONFIG = data;
            dfd.resolve(data);
        },
        (err) => {
            dfd.reject(`loadConfig file error${err.toString()}`);
        });
    } else {
        const fs = require('fs');
        fs.readFile(filename, { encoding: 'utf-8' }, (err, data) => {
            if (err) {
                dfd.reject(err);
            } else {
                data = parseConf(data, configFormat);
                global.CONFIG = data;
                dfd.resolve(data);
            }
        });
    }
    return dfd.promise;
};


module.exports = loadConfig;
