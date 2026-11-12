from app.modules.estimate.schema import EstimateRequest, EstimateResponse


class EstimateService:
    # Standard market reference price (IDR/Kg) and density (Kg/L)
    BENCHMARKS = {
        "PLASTIC": {"price": 4000.0, "density": 0.045, "item_weight": 0.025, "co2_factor": 1.5},
        "PAPER": {"price": 2500.0, "density": 0.12, "item_weight": 0.05, "co2_factor": 1.0},
        "METAL": {"price": 8500.0, "density": 0.25, "item_weight": 0.04, "co2_factor": 2.0},
        "GLASS": {"price": 1200.0, "density": 0.35, "item_weight": 0.25, "co2_factor": 0.5},
        "ORGANIC": {"price": 500.0, "density": 0.30, "item_weight": 0.1, "co2_factor": 0.2},
    }

    def calculate_estimate(self, req: EstimateRequest) -> EstimateResponse:
        cat = req.category.upper()
        meta = self.BENCHMARKS.get(cat, self.BENCHMARKS["PLASTIC"])

        weight = 0.0
        if req.volume_liters and req.volume_liters > 0:
            density = meta["density"] * (1.6 if req.is_compacted else 1.0)
            weight = req.volume_liters * density
        elif req.count_items and req.count_items > 0:
            weight = req.count_items * meta["item_weight"]
        else:
            weight = 1.0  # default baseline

        weight = round(weight, 2)
        total_value = round(weight * meta["price"], 2)
        carbon_offset = round(weight * meta["co2_factor"], 2)

        return EstimateResponse(
            category=cat,
            estimated_weight_kg=weight,
            unit_price_per_kg=meta["price"],
            estimated_total_value_idr=total_value,
            carbon_offset_kg=carbon_offset,
        )


estimate_service = EstimateService()
