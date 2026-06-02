import './style.css';
import {GeneratePlate, GetDefaults, RenderPlate, SaveCurrentPNG} from '../wailsjs/go/main/App';

const shapeModes = Object.freeze({
  circle: 'Circle',
  text: 'Text'
});
const statusText = Object.freeze({
  ready: 'Ready',
  starting: 'Starting',
  generating: 'Generating',
  rendering: 'Rendering',
  saving: 'Saving',
  enterText: 'Enter text first'
});
const defaultConfuserColor = '#569B9B';
const presets = ['#61B9B9', '#569B9B', '#929292', '#DDB85F'];
const state = {
  density: 2000,
  minDensity: 10,
  maxDensity: 7000,
  textDensity: 5500,
  shade: 170,
  minShade: 156,
  maxShade: 190,
  shapeMode: shapeModes.circle,
  shapeText: '',
  confuserColor: defaultConfuserColor,
  showOutline: false,
  image: '',
  placed: 0,
  requested: 2000,
  busy: false,
  status: statusText.ready
};

let renderTimer;
let generateTimer;

select('#app').innerHTML = `
  <div class="shell">
    <main class="workspace">
      <section class="plateWrap" aria-label="Generated plate">
        <div class="plateStage">
          <img id="plateImage" alt="Generated reverse colorblind plate">
          <div class="emptyState" id="emptyState">${statusText.generating}</div>
        </div>
      </section>
    </main>

    <aside class="controls">
      <section class="controlSection heroControls">
        <div>
          <p class="panelLabel">Plate</p>
          <h2 id="placedText">Placed 0 / 0</h2>
          <p class="statusLine" id="statusText">${statusText.ready}</p>
        </div>
        <button class="primary" id="generateBtn" type="button">Generate</button>
      </section>

      <section class="controlSection">
        <div class="sectionHeader">
          <h2>Hidden Shape</h2>
          <div class="segmented" role="radiogroup" aria-label="Hidden shape">
            <button class="segment active" id="circleMode" type="button" aria-pressed="true">${shapeModes.circle}</button>
            <button class="segment" id="textMode" type="button" aria-pressed="false">${shapeModes.text}</button>
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
          <span class="valueText" id="confuserValue">${defaultConfuserColor}</span>
        </div>
        <div class="swatches" id="swatches"></div>
        <div class="colorLine">
          <input id="confuserPicker" type="color" value="${defaultConfuserColor}" aria-label="Confuser color" title="Confuser color">
          <input id="confuserHex" type="text" value="${defaultConfuserColor}" maxlength="7" spellcheck="false">
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
`;

const els = {
  plateImage: select('#plateImage'),
  emptyState: select('#emptyState'),
  generateBtn: select('#generateBtn'),
  saveBtn: select('#saveBtn'),
  statusText: select('#statusText'),
  placedText: select('#placedText'),
  density: select('#density'),
  densityValue: select('#densityValue'),
  shade: select('#shade'),
  shadeValue: select('#shadeValue'),
  revealChip: select('#revealChip'),
  circleMode: select('#circleMode'),
  textMode: select('#textMode'),
  shapeText: select('#shapeText'),
  confuserPicker: select('#confuserPicker'),
  confuserHex: select('#confuserHex'),
  confuserValue: select('#confuserValue'),
  swatches: select('#swatches'),
  outlineToggle: select('#outlineToggle'),
  outlineState: select('#outlineState')
};

function select(selector) {
  return document.querySelector(selector);
}

function clampNumber(value, min, max) {
  return Math.min(Math.max(Number(value), min), max);
}

function isTextMode() {
  return state.shapeMode === shapeModes.text;
}

function hasShapeText() {
  return state.shapeText.trim() !== '';
}

function requestFromState() {
  return {
    density: state.density,
    shade: state.shade,
    shapeMode: state.shapeMode,
    shapeText: state.shapeText,
    confuserColor: state.confuserColor,
    showOutline: state.showOutline
  };
}

function setBusy(busy, status) {
  state.busy = busy;
  state.status = status;
  els.generateBtn.disabled = busy;
  els.saveBtn.disabled = busy || !state.image;
  els.statusText.textContent = status;
  document.body.classList.toggle('busy', busy);
}

function updateControls() {
  const reveal = `rgb(255, 0, ${state.shade})`;

  if (state.density > state.maxDensity) {
    state.density = state.maxDensity;
  }

  syncRangeAndNumber(els.density, els.densityValue, state.minDensity, state.maxDensity, state.density);
  syncRangeAndNumber(els.shade, els.shadeValue, state.minShade, state.maxShade, state.shade);
  updateShapeControls();
  updateColorControls(reveal);
  updatePlateStatus();
}

function syncRangeAndNumber(range, number, min, max, value) {
  range.min = min;
  range.max = max;
  range.value = value;
  number.min = min;
  number.max = max;
  if (document.activeElement !== number) {
    number.value = value;
  }
}

function updateShapeControls() {
  els.shapeText.value = state.shapeText;
  els.shapeText.disabled = !isTextMode();
  els.outlineState.textContent = state.showOutline ? 'Visible' : 'Hidden';
  els.outlineToggle.setAttribute('aria-pressed', String(state.showOutline));
  setSegment(els.circleMode, state.shapeMode === shapeModes.circle);
  setSegment(els.textMode, isTextMode());
}

function updateColorControls(reveal) {
  els.confuserPicker.value = state.confuserColor;
  if (document.activeElement !== els.confuserHex) {
    els.confuserHex.value = state.confuserColor;
  }
  els.confuserValue.textContent = state.confuserColor;
  els.revealChip.style.backgroundColor = reveal;
  document.documentElement.style.setProperty('--reveal-color', reveal);
  document.documentElement.style.setProperty('--confuser-color', state.confuserColor);
}

function updatePlateStatus() {
  els.placedText.textContent = `Placed ${state.placed} / ${state.requested}`;
}

function setSegment(button, active) {
  button.classList.toggle('active', active);
  button.setAttribute('aria-pressed', String(active));
}

function setPlate(result) {
  state.image = result.image;
  state.placed = result.placed;
  state.requested = result.requested;
  els.plateImage.src = result.image;
  els.plateImage.classList.add('loaded');
  els.emptyState.hidden = true;
  updateControls();
}

function isHex(value) {
  return /^#[0-9a-f]{6}$/i.test(value);
}

function normalizeHex(value) {
  const raw = value.trim();
  const prefixed = raw.startsWith('#') ? raw : `#${raw}`;
  return prefixed.toUpperCase();
}

async function generate() {
  if (isTextMode() && !hasShapeText()) {
    setBusy(false, statusText.enterText);
    return;
  }

  setBusy(true, statusText.generating);
  try {
    const result = await GeneratePlate(requestFromState());
    setPlate(result);
    setBusy(false, result.message);
  } catch (err) {
    setBusy(false, String(err));
  }
}

async function render() {
  setBusy(true, statusText.rendering);
  try {
    const result = await RenderPlate(requestFromState());
    setPlate(result);
    setBusy(false, result.message);
  } catch (err) {
    setBusy(false, String(err));
  }
}

function scheduleRender() {
  window.clearTimeout(renderTimer);
  renderTimer = window.setTimeout(render, 80);
}

function scheduleGenerate() {
  window.clearTimeout(generateTimer);
  generateTimer = window.setTimeout(generate, 450);
}

function setShapeMode(mode) {
  const previousMode = state.shapeMode;
  state.shapeMode = mode;
  if (isTextMode() && previousMode !== shapeModes.text && state.density < state.textDensity) {
    state.density = state.textDensity;
  }
  updateControls();
  if (isTextMode()) {
    els.shapeText.focus();
    if (hasShapeText()) {
      scheduleGenerate();
    }
  } else {
    scheduleGenerate();
  }
}

function setConfuserColor(value) {
  const hex = normalizeHex(value);
  if (!isHex(hex)) {
    return;
  }
  state.confuserColor = hex;
  updateControls();
  updateSwatches();
  scheduleRender();
}

function updateSwatches() {
  for (const button of els.swatches.querySelectorAll('.swatch')) {
    button.classList.toggle('active', button.dataset.color === state.confuserColor);
  }
}

function buildSwatches() {
  els.swatches.innerHTML = '';
  for (const preset of presets) {
    const button = document.createElement('button');
    button.type = 'button';
    button.className = 'swatch';
    button.style.backgroundColor = preset;
    button.dataset.color = preset;
    button.setAttribute('aria-label', preset);
    button.addEventListener('click', () => setConfuserColor(preset));
    els.swatches.append(button);
  }
  updateSwatches();
}

els.generateBtn.addEventListener('click', generate);
els.saveBtn.addEventListener('click', async () => {
  setBusy(true, statusText.saving);
  try {
    const path = await SaveCurrentPNG(requestFromState());
    setBusy(false, path ? `Saved ${path}` : 'Save cancelled');
  } catch (err) {
    setBusy(false, String(err));
  }
});

els.density.addEventListener('input', event => {
  state.density = Number(event.target.value);
  updateControls();
});

els.densityValue.addEventListener('change', event => {
  state.density = clampNumber(event.target.value, state.minDensity, state.maxDensity);
  updateControls();
});

els.shade.addEventListener('input', event => {
  state.shade = Number(event.target.value);
  updateControls();
  scheduleRender();
});

els.shadeValue.addEventListener('change', event => {
  state.shade = clampNumber(event.target.value, state.minShade, state.maxShade);
  updateControls();
  scheduleRender();
});

els.circleMode.addEventListener('click', () => setShapeMode(shapeModes.circle));
els.textMode.addEventListener('click', () => setShapeMode(shapeModes.text));

els.shapeText.addEventListener('input', event => {
  state.shapeText = event.target.value;
  if (!isTextMode()) {
    return;
  }
  if (!hasShapeText()) {
    setBusy(false, statusText.enterText);
    return;
  }
  scheduleGenerate();
});

els.confuserPicker.addEventListener('input', event => setConfuserColor(event.target.value));
els.confuserHex.addEventListener('change', event => setConfuserColor(event.target.value));
els.confuserHex.addEventListener('input', event => {
  const hex = normalizeHex(event.target.value);
  if (isHex(hex)) {
    state.confuserColor = hex;
    els.confuserPicker.value = hex;
    els.confuserValue.textContent = hex;
    updateSwatches();
    scheduleRender();
  }
});

els.outlineToggle.addEventListener('click', () => {
  state.showOutline = !state.showOutline;
  updateControls();
  scheduleRender();
});

async function boot() {
  buildSwatches();
  setBusy(true, statusText.starting);
  try {
    const defaults = await GetDefaults();
    Object.assign(state, defaults);
    updateControls();
    await generate();
  } catch (err) {
    setBusy(false, String(err));
  }
}

boot();
