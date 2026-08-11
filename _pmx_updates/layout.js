import './globals.css';

export const metadata = {
  title: 'Picomatch Go Engineering Foundry',
  description: 'Production validation platform and interactive match dashboard for picomatch-go',
};

export default function RootLayout({ children }) {
  return (
    <html lang="en">
      <body>
        {children}
      </body>
    </html>
  );
}
