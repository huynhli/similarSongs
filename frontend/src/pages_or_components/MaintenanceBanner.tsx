import React, { useEffect, useState } from 'react';
import $ from 'jquery';

interface MaintenanceData {
    is_active: boolean;
    message?: string;
    severity?: string;
}

export default function MaintenanceBanner() {
    const [maintenance, setMaintenance] = useState<MaintenanceData | null>(null);

    useEffect(() => {
        // Query the Go Fiber backend endpoint using jQuery
        $.getJSON('/api/cmsMaintenance')
        .done((data: MaintenanceData) => {
            if (data) {
            setMaintenance(data);
            }
        })
        .fail((_, textStatus, errorThrown) => {
            console.error("Failed to fetch maintenance status:", textStatus, errorThrown);
        });
    }, []);

    // Do not render anything if maintenance is inactive or data hasn't loaded
    if (!maintenance) {
        return null;
    }

    return (
        <div style={{ backgroundColor: '#ff4d4f', color: '#fff', padding: '12px', textAlign: 'center', fontWeight: 'bold' }}>
            <span>System Maintenance Notice: </span>
            <span>
                <div
                    dangerouslySetInnerHTML={{
                        __html:
                            maintenance.message ||
                            "<p>The system is currently undergoing scheduled maintenance.</p>",
                    }}
                />
            </span>
        </div>
    );
}