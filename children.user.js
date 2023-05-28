// ==UserScript==
// @name         wikitree children
// @namespace    https://github.com/yulvil/wikitree
// @version      0.2
// @description  Put siblings/children on separate lines
// @author       yulvil
// @match        http*://www.wikitree.com/wiki/*
// @match        http*://www.wikitree.com/index.php*
// @grant        none
// @updateURL    https://raw.githubusercontent.com/yulvil/wikitree/master/children.user.js
// ==/UserScript==

(function() {
    'use strict';

    var z = document.querySelectorAll("[itemprop=children]");

    for (let i=0; i<z.length; i++) {
        z[i].prepend(`${i+1}. `);
        z[i].prepend(document.createElement("div"));
    }

    z = document.querySelectorAll("[itemprop=sibling]");

    for (let i=0; i<z.length; i++) {
        z[i].prepend(`${i+1}. `);
        z[i].prepend(document.createElement("div"));
    }

})();
