# ColorBlind

A Go desktop application that generates Ishihara-style reverse colorblindness test plates.

The app fills a canvas with thousands of same-colored circles, then recolors circles inside a hidden shape with a slightly different shade of magenta. To someone with normal color vision, the plate looks uniform. To someone with a specific type of color blindness, the hidden shape becomes clearly visible.

## How it works

Based on [this theory](https://www.reddit.com/r/ColorBlind/comments/2ds5u1/this_is_a_reversecolorblind_test_normal_color/cjsxha1):

(hit ctrl-f and press either 3 or 4 or 5)

```
12292191298792287899799826279226826216879
27623444335345433545335453345545545348229
29263341986998212627268726281262715442971
26925438718926278962289287162892783352796
27723351981892628296771928127129685348962
81775542289219262891891897982898925459267
12683358962896271298678926899719773437127
96194439961296129612926292962296825532199
21773547177921721772276822792127963349612
79125433554355454355333445453443554347217
89267892677298912161279672716721912179677
```

As viewed by a color regular, its just a block of numbers. As viewed by someone for whom 3s and 5s look like 4s, its

```
12292191298792287899799826279226826216879
27624444444444444444444444444444444448229
29264441986998212627268726281262714442971
26924448718926278962289287162892784442796
27724441981892628296771928127129684448962
81774442289219262891891897982898924449267
12684448962896271298678926899719774447127
96194449961296129612926292962296824442199
21774447177921721772276822792127964449612
79124444444444444444444444444444444447217
89267892677298912161279672716721912179677
```

a clearly identifiable shape.

The app applies this same principle visually — using two nearly identical magenta shades that only become distinguishable to people with specific color vision deficiencies.

## Building

Requires Go and Fyne system libraries on Linux:

```bash
sudo apt-get install libgl-dev libx11-dev libxrandr-dev libxxf86vm-dev libxi-dev libxcursor-dev libxinerama-dev
```

```bash
go build -o colorblind ./cmd/colorblind
./colorblind
```

## Controls

- **Generate** — creates a new plate with randomized circle positions and a hidden shape
- **Density slider** — controls how many circles to place (10–2000)
- **Shade slider** — adjusts the magenta shade used for the hidden pattern; move it until the difference is barely perceptible to you
