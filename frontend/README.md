# TimeLedger Frontend

Nuxt 3 + Tailwind CSS + LINE LIFF Frontend for TimeLedger.

## Tech Stack

- **Framework**: Nuxt 3 (SSR)
- **Styling**: Tailwind CSS
- **State Management**: Pinia
- **Auth**: LINE LIFF
- **Utilities**: VueUse, HeadlessUI
- **Export**: html2canvas

## Project Structure

```
frontend/
├── components/         # Vue components
├── composables/       # Reusable composition functions
├── pages/             # Nuxt pages (file-based routing)
├── stores/            # Pinia stores
├── types/             # TypeScript type definitions
├── assets/            # Static assets (CSS, images)
├── server/            # Nuxt server middleware
├── nuxt.config.ts     # Nuxt configuration
├── tailwind.config.js # Tailwind configuration
└── package.json       # Dependencies
```

## Getting Started

### Prerequisites

- Node.js 18+
- npm or yarn or pnpm

### Installation

```bash
cd frontend
npm install
```

### Environment Variables

Create `.env` file:

```env
API_BASE_URL=http://localhost:8080/api/v1
LIFF_ID=YOUR_LIFF_ID
```

### Development

```bash
npm run dev
```

Open [http://localhost:3000](http://localhost:3000)

### Build

```bash
npm run build
```

### Preview Production Build

```bash
npm run preview
```

## Features

### Teacher Mobile App

- 📅 Unified calendar view across all centers
- ➕ Personal event management
- 🔔 Real-time notifications
- 👤 Profile management
- 📤 Export schedule as beautiful image

### Center Admin Dashboard

- 📊 Drag & drop scheduling grid
- 🔍 Smart matching for substitutes
- ✅ Exception approval workflow
- 📋 Resource management (rooms, courses, teachers)

## Design System

### Colors

- **Primary**: `#6366F1` (Indigo 500)
- **Secondary**: `#A855F7` (Purple 500)
- **Success**: `#10B981` (Emerald 500)
- **Critical**: `#F43F5E` (Rose 500)
- **Warning**: `#F59E0B` (Amber 500)

### Typography

- **Body**: 'Outfit', sans-serif
- **Headings**: 'Inter', sans-serif

### Components

- **Glass Cards**: Backdrop blur with subtle borders
- **Buttons**: Gradient backgrounds with hover effects
- **Inputs**: Rounded corners with focus rings

## Pages

- `/` - Landing page with LINE login
- `/admin/login` - Admin login page
- `/teacher/dashboard` - Teacher calendar dashboard
- `/teacher/profile` - Teacher profile page

## Stores

- `auth` - Authentication state
- `teacher` - Teacher data and schedule
- `notification` - Notifications

## Composables

- `useApi` - API request handler

## Middleware

- `auth-teacher` - Protect teacher routes
- `auth-admin` - Protect admin routes

## API Integration

All API calls use the `useApi()` composable:

```typescript
const api = useApi()
const response = await api.get<ResponseType>('/endpoint')
```

## LIFF Integration

LINE LIFF is integrated for teacher authentication:

```typescript
const liff = await import('@line/liff')
await liff.init({ liffId: config.public.liffId })
const profile = await liff.getProfile()
```

## Browser Support

- Chrome (latest)
- Firefox (latest)
- Safari (latest)
- Edge (latest)
- Mobile browsers (iOS Safari, Chrome Mobile)

## License

Copyright (c) 2026 TimeLedger Team. All rights reserved.
