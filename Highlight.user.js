// ==UserScript==
// @name         Highlight
// @namespace    http://tampermonkey.net/
// @version      0.2
// @description  try to take over the world!
// @author       You
//// @match        http://localhost/*
// @match        http://10.39.245.249:9080/rmweb/screen.vio*
// @icon         https://www.google.com/s2/favicons?domain=undefined.localhost
// @grant        GM_xmlhttpRequest
// @connect      *
// ==/UserScript==

(function() {
    'use strict';

    var hphm = document.querySelector('#hphm_right');
    if (!hphm) {
        return
    }

    var plate = hphm.textContent;
    var backend_host = 'http://localhost:8080';
    var url = backend_host + '/query?p=' + plate;

    GM_xmlhttpRequest({
        method: "GET",
        url: url,
        responseType: 'json',
        onload: function(response) {
            console.log(response.response);
            let resp = response.response;
            if (resp.result) {
                hphm.style.backgroundColor = 'red';
                alert(resp.desc);
            }
        },
        onerror: function(response) {
            console.log(response);
            alert('获取货车审批结果失败。');
        }
    });
})();