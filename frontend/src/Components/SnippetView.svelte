<script>
    import { onMount, beforeUpdate, createEventDispatcher } from 'svelte';
    import {GetUUID, SaveSnippet, DeleteSnippet, RenderSnippet} from "../../wailsjs/go/main/App.js";

    export let selectedSnippet = null;
    let doWrap = false, inEditMode, timer, newUuid, currentUuid, snippet, title, tags, created, renderedContent;

    const dispatch = createEventDispatcher();

    onMount(() => {
        reset();
    })

    beforeUpdate(() => {
        if (typeof selectedSnippet?.id !== 'undefined' && selectedSnippet.id !== currentUuid) {
            reset();
        }
    })

    function reset() {
        GetUUID().then(data => newUuid = data);
        currentUuid = selectedSnippet?.id || '';
        snippet = selectedSnippet?.body || '';
        title = selectedSnippet?.title || '';
        tags = selectedSnippet?.tags || [];
        created = selectedSnippet?.created || null;
        inEditMode = snippet === '';
        render();
    }

    function addNew() {
        selectedSnippet = null;
        reset();
        dispatch('message', { snippet: null });
    }

    function deleteSnippet() {
        if (currentUuid) {
            DeleteSnippet(currentUuid)
            addNew()
        }
    }

    const debounce = e => {
        clearTimeout(timer);
        timer = setTimeout(() => {
            snippet = e.target.value;
            saveSnippet();
        }, 550);
    }

    const saveSnippet = async () => {
        currentUuid = currentUuid || newUuid
        renderedContent = '';
        updateMeta();
        SaveSnippet(
            currentUuid,
            snippet.trim(),
            title.trim() !== '' ? title.trim() : 'Untitled',
            tags,
            created,
        ).then(data => {
            render();
            dispatch('message', { snippet: data.id });
        });
    }

    function updateMeta() {
        if (created === null) {
            created = (new Date()).toISOString();
        }
        const foundTitle = snippet.match(/^#\s+([\w ].*)/m);
        if (foundTitle !== null && foundTitle.length) {
            title = foundTitle[1].trim();
        }
        const foundFenced = [...snippet.matchAll(/^```(\w+)\n/gm)];
        if (foundFenced !== null && foundFenced.length) {
            let foundTags = [];
            foundFenced.forEach(function (found, index) {
                if (!foundTags.includes(found[1])) {
                    foundTags.push(found[1]);
                }
            })
            foundTags.sort();
            tags = foundTags;
        } else {
            tags = [];
        }
    }

    function render() {
        RenderSnippet(currentUuid).then(data => {
            renderedContent = data;
            const matches = document.querySelectorAll("div.rendered > pre > code > pre");
            matches.forEach((codeItem) => {
                codeItem.style.padding = '1rem';
            });
        });
    }
</script>

<section>
    <nav>
        <button on:click={addNew}>
            <span>new</span>
        </button>
        {#if snippet}
            <button on:click={() => { inEditMode = !inEditMode }} class:on={inEditMode}>
                <span>edit</span>
            </button>
        {/if}
        {#if inEditMode}
            <button on:click={() => doWrap = !doWrap} class:on={doWrap}>
                <span>word-wrap</span>
            </button>
        {/if}
        {#if snippet}
            <button on:click={deleteSnippet} class="delete">
                <span>delete</span>
            </button>
        {/if}
    </nav>

    <textarea
        wrap={doWrap ? "soft" : "off"} spellcheck="true"
        id="snippet" name="snippet" rows="5" cols="33"
        placeholder={'# Some optional title\n\n```php\n$some = "content";\n```'}
        on:keydown={debounce}
        class:hidden={!inEditMode}
    >{snippet}</textarea>

    <div class="rendered" class:hidden={inEditMode}>{@html renderedContent}</div>
</section>

<style>
    section {
        flex: 1 1 100%;
        display: flex;
        flex-flow: column;
        overflow: hidden;
    }
    nav {
        flex: 0 1 auto;
        display: flex;
        flex-direction: row;
        justify-content: space-evenly;
        padding: 1rem 0;
        margin: 0;
    }
    textarea {
        padding: 1rem;
        flex: 1 1 auto;
        border: 0;
        outline: 0;
        resize: none;
        background-color: #3e3e3e;
        color: #D6D0A5;
        display: flex;
    }
    button {
        cursor: pointer;
        white-space: nowrap;
        top: 0;
        border-radius: 4px;
        position: relative;
        border: none;
        background: none;
        transition: 0.1s all ease;
        border-bottom: 1px solid #D6D0A5;
        color: #f5f5f5;
        background: #A28364;
        width: 100px;
        height: 25px;
        font-size: 10px;
        text-transform: uppercase;
    }
    button:hover {
        border-bottom: 3px solid #D6D0A5;
        top: 0;
        height: 25px;
    }
    button:active, button.on:active {
        top: 0;
        border: 1px solid #727272;
        background-color: #96A672;
        color: #48635A;
        box-shadow: inset 0px 0px 4px rgb(0 0 0 / 15%);
    }
    button.on {
        border-bottom: 1px solid #48635A;
        background: #96A672;
        font-weight: bold;
    }
    button.on:hover {
        border-bottom: 3px solid #48635A;
        top: 0;
        height: 25px;
    }
    button.delete {
        border-bottom: 1px solid #ad5959;
        color: #860d0d;
        background: #e19d9d;
    }
    button.delete:hover {
        border-bottom: 3px solid #ad5959;
        top: 0;
        height: 25px;
    }
    button.delete:active {
        top: 0;
        border: 1px solid #b02b2b;
        background-color: #c56969;
        box-shadow: inset 0px 0px 4px rgb(0 0 0 / 15%);
        color: #fff;
    }
</style>
