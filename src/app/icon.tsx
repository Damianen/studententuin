import { ImageResponse } from 'next/og'

export const size = {
    width: 32,
    height: 32,
}
export const contentType = 'image/png'

export default function Icon() {
    return new ImageResponse(
        (
            <div
                style={{
                    width: '100%',
                    height: '100%',
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'center',
                    background: 'transparent',
                }}>
                <svg
                    width="28"
                    height="28"
                    viewBox="0 0 24 24"
                    fill="none"
                    xmlns="http://www.w3.org/2000/svg">
                    <path
                        d="M12 2C10.5 4 9 6.5 9 9C9 9.5 9.1 10 9.3 10.5C8.5 10.2 7.8 10 7 10C5.3 10 4 11 4 12.5C4 13.8 4.9 14.8 6 15.2C5.4 15.4 5 15.9 5 16.5C5 17.3 5.7 18 6.5 18H11V22H13V18H17.5C18.3 18 19 17.3 19 16.5C19 15.9 18.6 15.4 18 15.2C19.1 14.8 20 13.8 20 12.5C20 11 18.7 10 17 10C16.2 10 15.5 10.2 14.7 10.5C14.9 10 15 9.5 15 9C15 6.5 13.5 4 12 2Z"
                        fill="#22c55e"
                    />
                </svg>
            </div>
        ),
        {
            ...size,
        }
    )
}
