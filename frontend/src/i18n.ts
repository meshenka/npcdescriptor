export type Locale = 'en' | 'fr';

export const translations = {
  en: {
    title: 'NPC Descriptors',
    countLabel: 'Count (1-10):',
    generateBtn: 'Generate Descriptors',
    loading: 'Loading...',
    copyBtn: 'Copy All',
    copied: 'Copied!',
    fetchError: 'Failed to fetch descriptors.',
    clipboardError: 'Clipboard API not available (requires HTTPS or modern browser).',
    copyError: 'Failed to copy to clipboard.',
    historyTitle: 'Recent Rolls'
  },
  fr: {
    title: 'Descripteurs de PNJ',
    countLabel: 'Nombre (1-10) :',
    generateBtn: 'Générer les descripteurs',
    loading: 'Chargement...',
    copyBtn: 'Tout copier',
    copied: 'Copié !',
    fetchError: 'Échec de la récupération des descripteurs.',
    clipboardError: 'API Presse-papiers non disponible (HTTPS or navigateur moderne requis).',
    copyError: 'Échec de la copie dans le presse-papiers.',
    historyTitle: 'Tirages récents'
  }
};
