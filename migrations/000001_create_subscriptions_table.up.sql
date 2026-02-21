CREATE TABLE IF NOT EXISTS subscriptions (
                                             id UUID PRIMARY KEY,
                                             service_name VARCHAR(255) NOT NULL,
    price INTEGER NOT NULL CHECK (price > 0),
    user_id UUID NOT NULL,
    start_date VARCHAR(7) NOT NULL, -- формат MM-YYYY
    end_date VARCHAR(7),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
                             );

CREATE INDEX idx_subscriptions_user_id ON subscriptions(user_id);
CREATE INDEX idx_subscriptions_service_name ON subscriptions(service_name);
CREATE INDEX idx_subscriptions_dates ON subscriptions(start_date, end_date);
CREATE INDEX idx_subscriptions_user_dates ON subscriptions(user_id, start_date, end_date);