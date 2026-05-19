/**
 * HOWL Chat — Tutorial Data Orchestrator
 *
 * Merges all tutorial section data into a single TUTORIALS_DATA object.
 *
 * Depends on:
 *   tutorials-data-chat.js
 *   tutorials-data-settings.js
 *   tutorials-data-model.js
 *   tutorials-data-content.js
 *   tutorials-data-advanced.js
 *   tutorials-data-app.js
 *
 * @module tutorials-data
 */

'use strict';

(function() {
    const sections = [];

    if (window.TUTORIALS_DATA_CHAT)     sections.push(window.TUTORIALS_DATA_CHAT);
    if (window.TUTORIALS_DATA_SETTINGS) sections.push(window.TUTORIALS_DATA_SETTINGS);
    if (window.TUTORIALS_DATA_MODEL)    sections.push(window.TUTORIALS_DATA_MODEL);
    if (window.TUTORIALS_DATA_CONTENT)  sections.push(window.TUTORIALS_DATA_CONTENT);
    if (window.TUTORIALS_DATA_ADVANCED) sections.push(window.TUTORIALS_DATA_ADVANCED);
    if (window.TUTORIALS_DATA_APP)      sections.push(window.TUTORIALS_DATA_APP);

    window.TUTORIALS_DATA = {
        sections: sections
    };
})();
