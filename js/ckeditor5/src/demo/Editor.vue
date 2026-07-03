<template>
  <ckeditor
    id="cke"
    v-if="editor"
    v-model="modelValue"
    :editor="editor"
    :config="config"
    @ready="onEditorReady"
  />
</template>

<script setup lang="ts">
import "./content.css";
import {computed} from "vue";
import {Ckeditor} from "@ckeditor/ckeditor5-vue";

import type {EventInfo} from "ckeditor5";
import {
  Alignment,
  Autoformat,
  AutoImage,
  AutoLink,
  Base64UploadAdapter,
  BlockQuote,
  BlockToolbar,
  Bold,
  Bookmark,
  ClassicEditor,
  CloudServices,
  Code,
  CodeBlock,
  Essentials,
  EventInfo,
  FindAndReplace,
  Font,
  GeneralHtmlSupport,
  Heading,
  Highlight,
  HorizontalLine,
  Image,
  ImageCaption,
  ImageInsert,
  ImageResize,
  ImageStyle,
  ImageToolbar,
  ImageUpload,
  Indent,
  IndentBlock,
  Italic,
  Link,
  LinkImage,
  List,
  ListProperties,
  MediaEmbed,
  Mention,
  PageBreak,
  Paragraph,
  PasteFromOffice,
  PictureEditing,
  RemoveFormat,
  SpecialCharacters,
  SpecialCharactersEssentials,
  Strikethrough,
  Style,
  Subscript,
  Superscript,
  Table,
  TableCaption,
  TableCellProperties,
  TableColumnResize,
  TableProperties,
  TableToolbar,
  TextPartLanguage,
  TextTransformation,
  TodoList,
  Underline,
  WordCount,
} from "ckeditor5";
import "ckeditor5/ckeditor5.css";
import {EMOJIS_ARRAY, REDUCED_MATERIAL_COLORS} from "./colors.js";

const modelValue = defineModel();

const editor = computed(() => {
  return ClassicEditor;
});

/**
 * Enrich the special characters plugin with emojis.
 */
function SpecialCharactersEmoji(editor) {
  if (!editor.plugins.get("SpecialCharacters")) {
    return;
  }

  // Make sure Emojis are last on the list.
  this.afterInit = function () {
    editor.plugins.get("SpecialCharacters").addItems("Emoji", EMOJIS_ARRAY);
  };
}

/**
 * CKEditor 5 requires a license key.
 *
 * The "GPL" license key used below only allows you to use the open-source features.
 * To use the premium features, replace it with your commercial license key.
 * If you don't have one, you can get a trial license key from https://portal.ckeditor.com/checkout?plan=free.
 */
const LICENSE_KEY = "GPL";

const config = computed(() => {
  return {
    plugins: [
      Alignment,
      Autoformat,
      AutoImage,
      AutoLink,
      BlockQuote,
      Bold,
      Bookmark,
      CloudServices,
      Code,
      CodeBlock,
      Essentials,
      FindAndReplace,
      Font,
      GeneralHtmlSupport,
      Heading,
      Highlight,
      HorizontalLine,
      Image,
      ImageCaption,
      ImageInsert,
      ImageResize,
      ImageStyle,
      ImageToolbar,
      ImageUpload,
      Base64UploadAdapter,
      Indent,
      IndentBlock,
      Italic,
      Link,
      LinkImage,
      List,
      ListProperties,
      MediaEmbed,
      Mention,
      PageBreak,
      Paragraph,
      PasteFromOffice,
      PictureEditing,
      RemoveFormat,
      SpecialCharacters,
      SpecialCharactersEmoji,
      SpecialCharactersEssentials,
      Strikethrough,
      Style,
      Subscript,
      Superscript,
      Table,
      TableCaption,
      TableCellProperties,
      TableColumnResize,
      TableProperties,
      TableToolbar,
      TextPartLanguage,
      TextTransformation,
      TodoList,
      Underline,
      WordCount,
      BlockToolbar,
    ],
    licenseKey: LICENSE_KEY,

    toolbar: {
      shouldNotGroupWhenFull: true,
      items: [
        // --- Document-wide tools ----------------------------------------------------------------------
        "undo",
        "redo",
        "|",
        "findAndReplace",
        "selectAll",
        "|",

        // --- "Insertables" ----------------------------------------------------------------------------

        "link",
        "bookmark",
        "insertImage",
        "uploadImage",
        "insertTable",
        "blockQuote",
        "mediaEmbed",
        "codeBlock",
        "pageBreak",
        "horizontalLine",
        "specialCharacters",
        "-",

        // --- Block-level formatting -------------------------------------------------------------------
        "heading",
        "style",
        "|",

        // --- Basic styles, font and inline formatting -------------------------------------------------------
        "bold",
        "italic",
        "underline",
        "strikethrough",
        "highlight",
        {
          label: "Basic styles",
          icon: "text",
          items: [
            "fontSize",
            "fontFamily",
            "fontColor",
            "fontBackgroundColor",
            "superscript",
            "subscript",
            "code",
          ],
        },
        "removeFormat",
        "|",

        // --- Text alignment ---------------------------------------------------------------------------
        "alignment",
        "|",

        // --- Lists and indentation --------------------------------------------------------------------
        "bulletedList",
        "numberedList",
        "multilevelList",
        "todoList",
        "|",
        "outdent",
        "indent",
      ],
    },
    fontFamily: {
      supportAllValues: true,
    },
    fontSize: {
      options: [10, 12, 14, "default", 18, 20, 22],
      supportAllValues: true,
    },
    fontColor: {
      columns: 12,
      colors: REDUCED_MATERIAL_COLORS,
    },
    fontBackgroundColor: {
      columns: 12,
      colors: REDUCED_MATERIAL_COLORS,
    },

    heading: {
      options: [
        {
          model: "paragraph",
          title: "Paragraph",
          class: "ck-heading_paragraph",
        },
        {
          model: "heading1",
          view: "h1",
          title: "Heading 1",
          class: "ck-heading_heading1",
        },
        {
          model: "heading2",
          view: "h2",
          title: "Heading 2",
          class: "ck-heading_heading2",
        },
        {
          model: "heading3",
          view: "h3",
          title: "Heading 3",
          class: "ck-heading_heading3",
        },
        {
          model: "heading4",
          view: "h4",
          title: "Heading 4",
          class: "ck-heading_heading4",
        },
        {
          model: "heading5",
          view: "h5",
          title: "Heading 5",
          class: "ck-heading_heading5",
        },
        {
          model: "heading6",
          view: "h6",
          title: "Heading 6",
          class: "ck-heading_heading6",
        },
      ],
    },
    htmlSupport: {
      allow: [
        // Enables all HTML features.
        {
          name: /.*/,
          attributes: true,
          classes: true,
          styles: true,
        },
      ],
      disallow: [
        {
          attributes: [
            { key: /^on(.*)/i, value: true },
            {
              key: /.*/,
              value: /(\b)(on\S+)(\s*)=|javascript:|(<\s*)(\/*)script/i,
            },
            { key: /.*/, value: /data:(?!image\/(png|jpeg|gif|webp))/i },
          ],
        },
        { name: "script" },
      ],
    },
    image: {
      resizeOptions: [
        {
          name: "resizeImage:original",
          label: "Default image width",
          value: null,
        },
        {
          name: "resizeImage:50",
          label: "50% page width",
          value: "50",
        },
        {
          name: "resizeImage:75",
          label: "75% page width",
          value: "75",
        },
      ],
      toolbar: [
        "imageTextAlternative",
        "toggleImageCaption",
        "|",
        "imageStyle:inline",
        "imageStyle:wrapText",
        "imageStyle:breakText",
        "|",
        "resizeImage",
      ],
      insert: {
        integrations: ["url"],
      },
    },
    list: {
      properties: {
        styles: true,
        startIndex: true,
        reversed: true,
      },
    },
    link: {
      decorators: {
        toggleDownloadable: {
          mode: "manual",
          label: "Downloadable",
          attributes: {
            download: "file",
          },
        },
      },
      addTargetToExternalLinks: true,
      defaultProtocol: "https://",
    },
    mention: {
      feeds: [
        {
          marker: "@",
          feed: [{ id: "@cflores", avatar: "m_1", name: "Charles Flores" }],
          minimumCharacters: 1,
          itemRenderer: customMentionUserItemRenderer,
        },
        {
          marker: "#",
          feed: ["#american"],
        },
      ],
    },
    placeholder: "Type or paste your content here!",
    style: {
      definitions: [
        /* {
          name: 'Title',
          element: 'h1',
          classes: ['document-title']
        },
        {
          name: 'Subtitle',
          element: 'h2',
          classes: ['document-subtitle']
        }, */
        {
          name: "Callout",
          element: "p",
          classes: ["callout"],
        },
        {
          name: "Side quote",
          element: "blockquote",
          classes: ["side-quote"],
        },
        {
          name: "Needs clarification",
          element: "span",
          classes: ["needs-clarification"],
        },
        {
          name: "Wide spacing",
          element: "span",
          classes: ["wide-spacing"],
        },
        {
          name: "Small caps",
          element: "span",
          classes: ["small-caps"],
        },
        {
          name: "Code (dark)",
          element: "pre",
          classes: ["stylish-code", "stylish-code-dark"],
        },
        {
          name: "Code (bright)",
          element: "pre",
          classes: ["stylish-code", "stylish-code-bright"],
        },
        {
          name: "Multlevel List",
          element: "ol",
          classes: ["multilevel-list"],
        },
      ],
    },
    table: {
      contentToolbar: [
        "tableColumn",
        "tableRow",
        "mergeTableCells",
        "tableProperties",
        "tableCellProperties",
        "toggleTableCaption",
      ],
    },
    menuBar: {
      isVisible: true,
    },
  };
});

/*
 * Customizes the way the list of user suggestions is displayed.
 * Each user has an @id, a name and an avatar.
 */
function customMentionUserItemRenderer(item) {
  const itemElement = document.createElement("span");
  const avatar = document.createElement("img");
  const userNameElement = document.createElement("span");
  const fullNameElement = document.createElement("span");

  itemElement.classList.add("mention__item");

  avatar.src = `/assets/images/avatars/${item.avatar}.jpg`;

  userNameElement.classList.add("mention__item__user-name");
  userNameElement.textContent = item.id;

  fullNameElement.classList.add("mention__item__full-name");
  fullNameElement.textContent = item.name;

  itemElement.appendChild(avatar);
  itemElement.appendChild(userNameElement);
  itemElement.appendChild(fullNameElement);

  return itemElement;
}

function onEditorReady(eventInfo: EventInfo) {
  const e: any = eventInfo;
  e.ui.view.element
    .querySelector(".ck.ck-editor__main")
    .appendChild(e.plugins.get("WordCount").wordCountContainer);
}
</script>
