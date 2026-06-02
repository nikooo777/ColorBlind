(function(){const a=document.createElement("link").relList;if(a&&a.supports&&a.supports("modulepreload"))return;for(const i of document.querySelectorAll('link[rel="modulepreload"]'))f(i);new MutationObserver(i=>{for(const r of i)if(r.type==="childList")for(const g of r.addedNodes)g.tagName==="LINK"&&g.rel==="modulepreload"&&f(g)}).observe(document,{childList:!0,subtree:!0});function u(i){const r={};return i.integrity&&(r.integrity=i.integrity),i.referrerpolicy&&(r.referrerPolicy=i.referrerpolicy),i.crossorigin==="use-credentials"?r.credentials="include":i.crossorigin==="anonymous"?r.credentials="omit":r.credentials="same-origin",r}function f(i){if(i.ep)return;i.ep=!0;const r=u(i);fetch(i.href,r)}})();function D(e){return window.go.main.App.GeneratePlate(e)}function V(){return window.go.main.App.GetDefaults()}function O(e){return window.go.main.App.RenderPlate(e)}function $(e){return window.go.main.App.SaveCurrentPNG(e)}const d=Object.freeze({circle:"Circle",text:"Text"}),l=Object.freeze({ready:"Ready",starting:"Starting",generating:"Generating",rendering:"Rendering",saving:"Saving",enterText:"Enter text first"}),m="#569B9B",G=["#61B9B9","#569B9B","#929292","#DDB85F"],t={density:2e3,minDensity:10,maxDensity:7e3,textDensity:5500,shade:170,minShade:156,maxShade:190,shapeMode:d.circle,shapeText:"",confuserColor:m,showOutline:!1,image:"",placed:0,requested:2e3,busy:!1,status:l.ready};let T,C;s("#app").innerHTML=`
  <div class="shell">
    <main class="workspace">
      <section class="plateWrap" aria-label="Generated plate">
        <div class="plateStage">
          <img id="plateImage" alt="Generated reverse colorblind plate">
          <div class="emptyState" id="emptyState">${l.generating}</div>
        </div>
      </section>
    </main>

    <aside class="controls">
      <section class="controlSection heroControls">
        <div>
          <p class="panelLabel">Plate</p>
          <h2 id="placedText">Placed 0 / 0</h2>
          <p class="statusLine" id="statusText">${l.ready}</p>
        </div>
        <button class="primary" id="generateBtn" type="button">Generate</button>
      </section>

      <section class="controlSection">
        <div class="sectionHeader">
          <h2>Hidden Shape</h2>
          <div class="segmented" role="radiogroup" aria-label="Hidden shape">
            <button class="segment active" id="circleMode" type="button" aria-pressed="true">${d.circle}</button>
            <button class="segment" id="textMode" type="button" aria-pressed="false">${d.text}</button>
          </div>
        </div>
        <input class="textInput" id="shapeText" type="text" maxlength="18" placeholder="Enter hidden text" disabled>
      </section>

      <section class="controlSection">
        <div class="labelRow">
          <label for="density">Density</label>
          <input class="numberInput" id="densityValue" type="number" min="10" max="2000" step="10" value="2000">
        </div>
        <input id="density" type="range" min="10" max="2000" step="10" value="2000">
      </section>

      <section class="controlSection">
        <div class="labelRow">
          <label for="shade">Reveal Shade</label>
          <input class="numberInput" id="shadeValue" type="number" min="156" max="190" step="1" value="170">
        </div>
        <div class="shadePreview" aria-hidden="true">
          <span class="mainMagenta"></span>
          <span class="revealMagenta" id="revealChip"></span>
        </div>
        <input id="shade" type="range" min="156" max="190" step="1" value="170">
      </section>

      <section class="controlSection">
        <div class="labelRow">
          <label for="confuserHex">Confuser</label>
          <span class="valueText" id="confuserValue">${m}</span>
        </div>
        <div class="swatches" id="swatches"></div>
        <div class="colorLine">
          <input id="confuserPicker" type="color" value="${m}" aria-label="Confuser color" title="Confuser color">
          <input id="confuserHex" type="text" value="${m}" maxlength="7" spellcheck="false">
        </div>
      </section>

      <section class="controlSection split compactToggle">
        <div>
          <h2>Show Shape</h2>
          <p id="outlineState">Hidden</p>
        </div>
        <button class="switch" id="outlineToggle" type="button" aria-pressed="false">
          <span></span>
        </button>
      </section>

      <div class="sidebarFooter">
        <button class="secondary" id="saveBtn" type="button">Save PNG</button>
      </div>
    </aside>
  </div>
`;const n={plateImage:s("#plateImage"),emptyState:s("#emptyState"),generateBtn:s("#generateBtn"),saveBtn:s("#saveBtn"),statusText:s("#statusText"),placedText:s("#placedText"),density:s("#density"),densityValue:s("#densityValue"),shade:s("#shade"),shadeValue:s("#shadeValue"),revealChip:s("#revealChip"),circleMode:s("#circleMode"),textMode:s("#textMode"),shapeText:s("#shapeText"),confuserPicker:s("#confuserPicker"),confuserHex:s("#confuserHex"),confuserValue:s("#confuserValue"),swatches:s("#swatches"),outlineToggle:s("#outlineToggle"),outlineState:s("#outlineState")};function s(e){return document.querySelector(e)}function P(e,a,u){return Math.min(Math.max(Number(e),a),u)}function p(){return t.shapeMode===d.text}function y(){return t.shapeText.trim()!==""}function x(){return{density:t.density,shade:t.shade,shapeMode:t.shapeMode,shapeText:t.shapeText,confuserColor:t.confuserColor,showOutline:t.showOutline}}function o(e,a){t.busy=e,t.status=a,n.generateBtn.disabled=e,n.saveBtn.disabled=e||!t.image,n.statusText.textContent=a,document.body.classList.toggle("busy",e)}function c(){const e=`rgb(255, 0, ${t.shade})`;t.density>t.maxDensity&&(t.density=t.maxDensity),L(n.density,n.densityValue,t.minDensity,t.maxDensity,t.density),L(n.shade,n.shadeValue,t.minShade,t.maxShade,t.shade),N(),A(e),R()}function L(e,a,u,f,i){e.min=u,e.max=f,e.value=i,a.min=u,a.max=f,document.activeElement!==a&&(a.value=i)}function N(){n.shapeText.value=t.shapeText,n.shapeText.disabled=!p(),n.outlineState.textContent=t.showOutline?"Visible":"Hidden",n.outlineToggle.setAttribute("aria-pressed",String(t.showOutline)),M(n.circleMode,t.shapeMode===d.circle),M(n.textMode,p())}function A(e){n.confuserPicker.value=t.confuserColor,document.activeElement!==n.confuserHex&&(n.confuserHex.value=t.confuserColor),n.confuserValue.textContent=t.confuserColor,n.revealChip.style.backgroundColor=e,document.documentElement.style.setProperty("--reveal-color",e),document.documentElement.style.setProperty("--confuser-color",t.confuserColor)}function R(){n.placedText.textContent=`Placed ${t.placed} / ${t.requested}`}function M(e,a){e.classList.toggle("active",a),e.setAttribute("aria-pressed",String(a))}function E(e){t.image=e.image,t.placed=e.placed,t.requested=e.requested,n.plateImage.src=e.image,n.plateImage.classList.add("loaded"),n.emptyState.hidden=!0,c()}function B(e){return/^#[0-9a-f]{6}$/i.test(e)}function H(e){const a=e.trim();return(a.startsWith("#")?a:`#${a}`).toUpperCase()}async function b(){if(p()&&!y()){o(!1,l.enterText);return}o(!0,l.generating);try{const e=await D(x());E(e),o(!1,e.message)}catch(e){o(!1,String(e))}}async function I(){o(!0,l.rendering);try{const e=await O(x());E(e),o(!1,e.message)}catch(e){o(!1,String(e))}}function h(){window.clearTimeout(T),T=window.setTimeout(I,80)}function v(){window.clearTimeout(C),C=window.setTimeout(b,450)}function k(e){const a=t.shapeMode;t.shapeMode=e,p()&&a!==d.text&&t.density<t.textDensity&&(t.density=t.textDensity),c(),p()?(n.shapeText.focus(),y()&&v()):v()}function S(e){const a=H(e);!B(a)||(t.confuserColor=a,c(),w(),h())}function w(){for(const e of n.swatches.querySelectorAll(".swatch"))e.classList.toggle("active",e.dataset.color===t.confuserColor)}function q(){n.swatches.innerHTML="";for(const e of G){const a=document.createElement("button");a.type="button",a.className="swatch",a.style.backgroundColor=e,a.dataset.color=e,a.setAttribute("aria-label",e),a.addEventListener("click",()=>S(e)),n.swatches.append(a)}w()}n.generateBtn.addEventListener("click",b);n.saveBtn.addEventListener("click",async()=>{o(!0,l.saving);try{const e=await $(x());o(!1,e?`Saved ${e}`:"Save cancelled")}catch(e){o(!1,String(e))}});n.density.addEventListener("input",e=>{t.density=Number(e.target.value),c()});n.densityValue.addEventListener("change",e=>{t.density=P(e.target.value,t.minDensity,t.maxDensity),c()});n.shade.addEventListener("input",e=>{t.shade=Number(e.target.value),c(),h()});n.shadeValue.addEventListener("change",e=>{t.shade=P(e.target.value,t.minShade,t.maxShade),c(),h()});n.circleMode.addEventListener("click",()=>k(d.circle));n.textMode.addEventListener("click",()=>k(d.text));n.shapeText.addEventListener("input",e=>{if(t.shapeText=e.target.value,!!p()){if(!y()){o(!1,l.enterText);return}v()}});n.confuserPicker.addEventListener("input",e=>S(e.target.value));n.confuserHex.addEventListener("change",e=>S(e.target.value));n.confuserHex.addEventListener("input",e=>{const a=H(e.target.value);B(a)&&(t.confuserColor=a,n.confuserPicker.value=a,n.confuserValue.textContent=a,w(),h())});n.outlineToggle.addEventListener("click",()=>{t.showOutline=!t.showOutline,c(),h()});async function F(){q(),o(!0,l.starting);try{const e=await V();Object.assign(t,e),c(),await b()}catch(e){o(!1,String(e))}}F();
