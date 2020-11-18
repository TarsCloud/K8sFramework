/* jshint indent: 2 */

module.exports = function(sequelize, DataTypes) {
  return sequelize.define('t_server_notifys', {
    id: {
      type: DataTypes.INTEGER(11),
      allowNull: false,
      primaryKey: true,
      autoIncrement: true
    },
    server_name: {
      type: DataTypes.STRING(50),
      allowNull: true
    },
    server_id: {
      type: DataTypes.STRING(100),
      allowNull: true
    },
    thread_id: {
      type: DataTypes.STRING(20),
      allowNull: true
    },
    command: {
      type: DataTypes.STRING(50),
      allowNull: true
    },
    result: {
      type: DataTypes.TEXT,
      allowNull: true
    },
    notifytime: {
      type: DataTypes.DATE,
      allowNull: true
    }
  }, {
    tableName: 't_server_notifys'
  });
};
