<script lang="ts">
  import { Footer, FooterCopyright, FooterLinkGroup, FooterLink, DarkMode } from "flowbite-svelte";
</script>

<Footer>
  <FooterCopyright href="/" by="Chocolingo™" year={2026} />
  <FooterLinkGroup class="mt-3 flex flex-wrap items-center text-sm text-gray-500 sm:mt-0 dark:text-gray-400">
    <FooterLink href="/">About</FooterLink>
    <FooterLink href="/">Privacy Policy</FooterLink>
    <FooterLink href="/">Licensing</FooterLink>
    <FooterLink href="/">Contact</FooterLink>
    <DarkMode />
  </FooterLinkGroup>
</Footer>