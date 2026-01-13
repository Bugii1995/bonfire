<button class="neobrutal-btn">
  <slot />
</button>

<style>
.neobrutal-btn {
  background: #ffd60a;
  border: 3px solid black;
  box-shadow: 4px 4px 0 black;
  padding: 12px 20px;
  font-weight: bold;
  cursor: pointer;
  transition:
  transform 0.15s cubic-bezier(0.22, 1, 0.36, 1),
  box-shadow 0.15s cubic-bezier(0.22, 1, 0.36, 1),
  background-color 0.15s ease;
}

.neobrutal-btn:hover {
  transform: translate(-2px, -2px);
  box-shadow: 6px 6px 0 black;
}

.neobrutal-btn:active {
  transform: translate(2px, 2px);
  box-shadow: 2px 2px 0 black;
}
</style>
