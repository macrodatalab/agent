var AWS = require('aws-sdk');

// get reference to S3 client
var s3 = new AWS.S3();

var Promise = require('promise');
var GoogleSpreadsheet = require("google-spreadsheet");

exports.handler = function(event, context) {
    var my_sheet = new GoogleSpreadsheet('1h41KqZubi0OfwwOKX0sg55RLQrtxps-LyqASaWPP9jE')

    var authFlow = Promise.denodeify(my_sheet.setAuth);
    var getInfo = Promise.denodeify(my_sheet.getInfo);

    s3.getObject({Bucket: 'macrodatalab-secret', Key: 'account'}, function (err, data) {
        cred = JSON.parse(data.Body.toString());
        authFlow(cred.email, cred.secret)
            .then(function (err) {
                return getInfo()
            })
            .then(function (sheet_info) {
                console.log(sheet_info.title + ' is loaded');
                var sheet1 = sheet_info.worksheets[0];
                return sheet1;
            })
            .then(function (one_sheet) {
                console.log('begin update...');
                var addRow = Promise.denodeify(one_sheet.addRow);
                var result = [];
                for (var i = 0; i < event.length; i++) {
                    console.log(event[i]);
                    result[i] = addRow(event[i]);
                }
                return Promise.all(result);
            })
            .done(function (junk) {
                context.done();
            });
    });
};
