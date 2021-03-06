/**
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

const _ = require('lodash');

module.exports = function (sequelize, DataTypes) {
  return sequelize.define('t_apply_cache_module_conf', {
    id: {
      type: DataTypes.INTEGER(11),
      allowNull: false,
      primaryKey: true,
      autoIncrement: true,
    },
    module_id: {
      type: DataTypes.INTEGER(11),
      allowNull: false,
    },
    apply_id: {
      type: DataTypes.INTEGER(11),
      allowNull: false,
      unique: 'applyModule',
    },
    module_name: {
      type: DataTypes.STRING(100),
      allowNull: false,
      unique: 'applyModule',
      defaultValue: '',
    },
    status: {
      type: DataTypes.INTEGER(4),
      allowNull: false,
      defaultValue: 0,
    },
    area: {
      type: DataTypes.STRING(50),
      allowNull: false,
      defaultValue: '',
    },
    idc_area: {
      type: DataTypes.STRING(50),
      allowNull: false,
      defaultValue: '',
    },
    set_area: {
      type: DataTypes.STRING(100),
      allowNull: false,
      defaultValue: '',
      get() {
        return _.compact(this.getDataValue('set_area') ? this.getDataValue('set_area').split(',') : []);
      },
      set(val) {
        return this.setDataValue('set_area', val.join(','));
      },
    },
    admin: {
      type: DataTypes.STRING(255),
      allowNull: false,
      defaultValue: '',
    },
    // 应用场景 1、 cache， 2、cache 加数据库
    cache_module_type: {
      type: DataTypes.INTEGER(4),
      allowNull: true,
      defaultValue: 0,
    },
    per_record_avg: {
      type: DataTypes.INTEGER(11),
      allowNull: false,
      defaultValue: 0,
    },
    total_record: {
      type: DataTypes.INTEGER(11),
      allowNull: false,
      defaultValue: 0,
    },
    max_read_flow: {
      type: DataTypes.INTEGER(11),
      allowNull: false,
      defaultValue: 0,
    },
    max_write_flow: {
      type: DataTypes.INTEGER(11),
      allowNull: false,
      defaultValue: 0,
    },
    // 1、内存数据库不淘汰， 2、临时缓存过期淘汰，3、热点缓存容量淘汰
    key_type: {
      type: DataTypes.INTEGER(4),
      allowNull: true,
      defaultValue: 0,
    },
    dbAccessServant: {
      type: DataTypes.STRING(150),
      allowNull: false,
      defaultValue: '',
    },
    module_remark: {
      type: DataTypes.TEXT,
      allowNull: true,
    },
    create_person: {
      type: DataTypes.STRING(50),
      allowNull: false,
      defaultValue: '',
    },
    modify_time: {
      type: DataTypes.DATE,
      defaultValue: DataTypes.NOW,
    },
  }, {
    tableName: 't_apply_cache_module_conf',
    timestamps: false,
  });
};
