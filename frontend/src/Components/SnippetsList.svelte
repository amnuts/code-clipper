<script>
    import {beforeUpdate, createEventDispatcher} from 'svelte';

    export let snippets = [], totalSnippets = 0, selectedSnippetId = '', selectedTag = '';
    let currentTag = '', filteredSnippets = [];
    const dispatch = createEventDispatcher();

    beforeUpdate(() => {
        snippets = Object.keys(snippets)
            .sort((keyA, keyB) => {
                const updatedA = Date.parse(snippets[keyA]['updated']);
                const updatedB = Date.parse(snippets[keyB]['updated']);
                if (updatedA > updatedB) {
                    return -1;
                }
                if (updatedA < updatedB) {
                    return 1;
                }
                return 0;
            })
            .reduce((obj, key, cIndex) => {
                obj[key] = snippets[key];
                return obj;
            }, {});

        if (selectedTag !== currentTag) {
            if (!selectedTag) {
                filteredSnippets = snippets;
                return;
            }
            filteredSnippets = Object.keys(snippets)
                .filter((key) => {
                    return snippets[key]['tags'].includes(selectedTag);
                })
                .reduce((obj, key, cIndex) => {
                    obj[key] = snippets[key];
                    return obj;
                }, {});
            if (!filteredSnippets.hasOwnProperty(selectedSnippetId)) {
                selectedSnippetId = '';
                dispatch('message', { snippet: '' });
            }
            currentTag = selectedTag;
        }
    })

    function selectSnippet(id) {
        selectedSnippetId = id;
        dispatch('message', { snippet: id });
    }
</script>

{#if totalSnippets}
<section>
{#each Object.entries(filteredSnippets) as [id, snippet]}
    <div on:click={() => selectSnippet(id)} class:selected={selectedSnippetId === id}>
        <p>{snippet.title}</p>
        <time>{new Date(snippet.updated).toLocaleString()}</time>
        {#if snippet.tags.length}
            <ol>
                {#each snippet.tags as tag, t }
                    <li>{tag}</li>
                {/each}
            </ol>
        {/if}
    </div>
{/each}
</section>
{/if}

<style>
    section {
        flex: 1 1 auto;
        min-width: 250px;
        padding: 1rem;
        border-right: 1px solid rgba(72, 99, 90, 30%);
        overflow-y: auto;
    }
    div {
        display: flex;
        flex-direction: column;
        margin-bottom: 0.75rem;
        padding-bottom: 0.25rem;
        cursor: pointer;
        border-bottom: 1px solid transparent;
    }
    div:hover {
        color: #D6D0A5;
        border-bottom: 1px dashed rgba(214, 208, 165, 50%);
    }
    div:last-of-type {
        margin-bottom: 0;
    }
    p, time, ol {
        display: inline-flex;
        margin: 0;
    }
    time, li {
        font-size: 80%;
        color: #A28364;
    }
    ol {
        list-style: none;
        margin: 0;
        padding: 0;
    }
    li {
        display: inline-block;
        color: #A28364;
        margin-left: 1rem;
    }
    li:first-of-type {
        margin-left: 0;
    }
    .selected {
        color: #96A672 !important;
        border-bottom: 1px solid #96A672 !important;
    }
</style>
