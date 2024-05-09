# {{ .header }}

{{ env "GREETING" }}, {{ if eq (env "NAME") "" }}friend{{ else }}{{ env "NAME" }}{{ end }}!

## Section 1

Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Laoreet suspendisse interdum consectetur libero id faucibus. Nec feugiat in fermentum posuere urna. Ut eu sem integer vitae justo eget. Feugiat scelerisque varius morbi enim nunc faucibus a pellentesque sit. In massa tempor nec feugiat nisl pretium fusce id velit. Dolor sit amet consectetur adipiscing elit ut aliquam. Venenatis cras sed felis eget velit aliquet sagittis id. Felis bibendum ut tristique et egestas. Nibh tortor id aliquet lectus. Sodales ut eu sem integer.

### Section 1.1

- bullet point 1
- bullet point 2
  - sub bullet point

--- 

### Section 1.2

1. numbered point 1
2. numbered point 2
   - sub bullet point

## Section 2

Inline code snippet: `pushd /tmp && ls -l && popd`

```yaml
# Multiline code snippet
key1: value
key2:
    - item1
    - item2
key3: 3
key4:
    subkey1: ["subvalue1"]
    subkey2: subvalue2
```

```bash
# Another multiline code snippet
echo "Hello, world!"
```

## Section 3

> Montes nascetur ridiculus mus mauris. Adipiscing bibendum est ultricies integer quis auctor elit. Morbi blandit cursus risus at ultrices mi tempus imperdiet. Eget dolor morbi non arcu risus. Interdum velit euismod in pellentesque massa placerat. Et magnis dis parturient montes nascetur. Blandit massa enim nec dui nunc mattis enim. At ultrices mi tempus imperdiet nulla malesuada pellentesque elit eget. Arcu cursus euismod quis viverra. Enim ut tellus elementum sagittis vitae. Nulla facilisi nullam vehicula ipsum. Curabitur gravida arcu ac tortor dignissim. Feugiat pretium nibh ipsum consequat.
> 
> - Author

**Duis at tellus at urna condimentum mattis pellentesque id.** Gravida quis blandit turpis cursus in hac. _Dui id ornare arcu odio ut sem nulla pharetra._ ~~Libero enim sed faucibus turpis in eu mi.~~

## Section 4

| Header 1 | Header 2 | Header 3 |
|----------|----------|----------|
| Value 1  | Value 2  | Value 3  |
| Value 4  | Value 5  | Value 6  |

## Section 5

[Link to Google](https://www.google.com)

![Image](https://via.placeholder.com/150)

## Section 6

- [ ] Task 1
- [x] Task 2
- [ ] Task 3
